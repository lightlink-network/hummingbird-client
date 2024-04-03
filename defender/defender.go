package defender

import (
	"fmt"
	"hummingbird/node"
	"hummingbird/utils"
	"math/big"
	"strings"
	"time"

	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"hummingbird/node/contracts"
	chainOracleContract "hummingbird/node/contracts/ChainOracle.sol"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
)

const (
	ErrNoDataCommitment = "no data commitment has been generated for the provided height"
)

type Opts struct {
	Logger      *slog.Logger
	WorkerDelay time.Duration
}

type Defender struct {
	*node.Node
	Opts *Opts
}

func NewDefender(node *node.Node, opts *Opts) *Defender {
	return &Defender{Node: node, Opts: opts}
}

// Start starts the defender.
func (d *Defender) Start() error {
	err := d.startDefender()
	return err
}

// Starts the main defender loop.
func (d *Defender) startDefender() error {
	ticker := time.NewTicker(d.Opts.WorkerDelay)
	defer ticker.Stop()

	for range ticker.C {
		scanRanges := d.Ethereum.GetChallengeWindowBlockRanges()

		for i := 0; i < len(scanRanges); i++ {
			if len(scanRanges[i]) != 2 {
				return fmt.Errorf("invalid block range")
			}
			daChallenges, err := d.getDAChallenges(scanRanges[i][0], scanRanges[i][1], contracts.ChallengeDAStatusChallengerInitiated)
			if err != nil {
				return fmt.Errorf("error getting DA challenges: %w", err)
			}
			d.defendDAChallenges(*daChallenges)

			l2HeaderChallenges, err := d.getL2HeaderChallenges(scanRanges[i][0], scanRanges[i][1], contracts.ChallengeL2HeaderStatusChallengerInitiated)
			if err != nil {
				return fmt.Errorf("error getting L2 header challenges: %w", err)
			}
			d.defendL2HeaderChallenges(*l2HeaderChallenges)
		}
	}
	return nil
}

// Gets DA challenge events from Challenge.sol for the given block range and status.
func (d *Defender) getDAChallenges(startblock, endblock uint64, status uint8) (*challengeContract.ChallengeChallengeDAUpdateIterator, error) {
	log := d.Opts.Logger.With(
		"startblock", startblock,
		"endblock", endblock,
		"range", endblock-startblock,
		"status", status,
	)
	log.Debug("Starting log scan for historic pending DA challenges")

	opts := &bind.FilterOpts{
		Start: startblock,
		End:   &endblock,
	}

	challenges, err := d.Ethereum.FilterChallengeDAUpdate(opts, nil, nil, []uint8{status})
	if err != nil {
		return nil, err
	}
	defer challenges.Close()

	log.Debug("Finished log scan for historic pending DA challenges")
	return challenges, nil
}

// Defends multiple DA challenge events by iterating through the given iterator and attempting
// to defend each challenge.
func (d *Defender) defendDAChallenges(c challengeContract.ChallengeChallengeDAUpdateIterator) {
	for c.Next() {
		err := d.defendDAChallenge(*c.Event)
		if err != nil {
			d.Opts.Logger.Error("error defending DA challenge", "error", err)
		}
	}
}

// Defends a DA challenge event.
func (d *Defender) defendDAChallenge(c challengeContract.ChallengeChallengeDAUpdate) error {
	// ensure the challenge is in the correct status to be defended
	challengeInfo, err := d.Ethereum.GetDataRootInclusionChallenge(c.BlockHash, uint8(c.PointerIndex.Uint64()))
	if err != nil {
		return fmt.Errorf("error getting data root inclusion challenge: %w", err)
	}
	if challengeInfo.Status != contracts.ChallengeDAStatusChallengerInitiated {
		return fmt.Errorf("challenge is not pending")
	}

	blockHash := common.BytesToHash(c.BlockHash[:])
	statusString := contracts.DAChallengeStatusToString(c.Status)

	log := d.Opts.Logger.With(
		"blockHash", blockHash.Hex(),
		"blockIndex", c.BlockIndex,
		"expiry", time.Unix(c.Expiry.Int64(), 0).Format(time.RFC1123Z),
		"statusEnum", c.Status,
		"statusString", statusString,
	)
	log.Info("Attempting to defend pending DA challenge")

	// attempt to defend the challenge by submitting a tx to the Challenge contract
	tx, err := d.DefendDA(c.BlockHash, uint8(c.PointerIndex.Uint64()))
	if err != nil {
		if strings.Contains(err.Error(), ErrNoDataCommitment) {
			log.Info("Pending DA challenge is awaiting data commitment from Celestia validators, will retry later")
			return nil
		} else {
			return fmt.Errorf("error defending DA challenge: %w", err)
		}
	}

	log.Info("Pending DA challenge defended successfully", "tx", tx.Hash().Hex())

	return nil
}

// Attempts to defend a DA challenge for the given block hash.
//
// Queries Celestia for a proof of data availability and submits a tx to the Challenge contract.
func (d *Defender) DefendDA(block common.Hash, pointerIndex uint8) (*types.Transaction, error) {
	proof, err := d.GetDAProof(block, pointerIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to prove data availability: %w", err)
	}
	key, err := d.Ethereum.DataRootInclusionChallengeKey(nil, block, pointerIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get data root inclusion challenge key: %w", err)
	}
	return d.Ethereum.DefendDataRootInclusion(key, *proof)
}

// Gets the Celestia pointer for the given block hash and queries Celestia for a proof
// of data availability.
func (d *Defender) GetDAProof(block common.Hash, pointerIndex uint8) (*challengeContract.ChallengeDataAvailabilityChallengeDAProof, error) {
	pointers, err := d.GetDAPointer(block)
	if err != nil {
		return nil, fmt.Errorf("failed to get Celestia pointer: %w", err)
	}
	if pointers == nil {
		return nil, fmt.Errorf("no Celestia pointer found")
	}
	commit, err := d.Ethereum.GetBlobstreamCommitment(int64(pointers[pointerIndex].Height))
	if err != nil {
		return nil, fmt.Errorf("failed to get blobstream commitment: %w", err)
	}
	proof, err := d.Celestia.GetProof(pointers[pointerIndex], commit.StartBlock, commit.EndBlock, *commit.ProofNonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	p := &challengeContract.ChallengeDataAvailabilityChallengeDAProof{
		RootNonce: commit.ProofNonce,
		DataRootTuple: challengeContract.DataRootTuple{
			Height:   big.NewInt(int64(pointers[pointerIndex].Height)),
			DataRoot: proof.Tuple.DataRoot,
		},
		Proof: challengeContract.BinaryMerkleProof{
			SideNodes: proof.WrappedProof.SideNodes,
			Key:       proof.WrappedProof.Key,
			NumLeaves: proof.WrappedProof.NumLeaves,
		},
	}
	return p, nil
}

// Gets L2 Header challenge events from Challenge.sol for the given block range and status.
func (d *Defender) getL2HeaderChallenges(startblock, endblock uint64, status uint8) (*challengeContract.ChallengeL2HeaderChallengeUpdateIterator, error) {
	log := d.Opts.Logger.With(
		"startblock", startblock,
		"endblock", endblock,
		"range", endblock-startblock,
		"status", status,
	)
	log.Debug("Starting log scan for historic pending L2 header challenges")

	opts := &bind.FilterOpts{
		Start: startblock,
		End:   &endblock,
	}

	challenges, err := d.Ethereum.FilterL2HeaderChallengeUpdate(opts, nil, nil, []uint8{status})
	if err != nil {
		return nil, err
	}
	defer challenges.Close()

	log.Debug("Finished log scan for historic pending L2 header challenges")
	return challenges, nil
}

// Defends multiple L2 header challenge events by iterating through the given iterator and attempting to
// defend each challenge.
func (d *Defender) defendL2HeaderChallenges(c challengeContract.ChallengeL2HeaderChallengeUpdateIterator) {
	for c.Next() {
		err := d.defendL2HeaderChallenge(*c.Event)
		if err != nil {
			d.Opts.Logger.Error("error defending L2 header challenge", "error", err)
		}
	}
}

// Defends an L2 header challenge event.
func (d *Defender) defendL2HeaderChallenge(c challengeContract.ChallengeL2HeaderChallengeUpdate) error {
	// ensure the challenge is in the correct status to be defended
	challengeInfo, err := d.Ethereum.GetL2HeaderChallenge(c.ChallengeHash)
	if err != nil {
		return fmt.Errorf("error getting L2 header challenge: %w", err)
	}
	if challengeInfo.Status != contracts.ChallengeL2HeaderStatusChallengerInitiated {
		return fmt.Errorf("challenge is not pending")
	}

	rblock := common.BytesToHash(c.Rblock[:])
	l2BlockNum := c.L2Number

	log := d.Opts.Logger.With(
		"rblock", rblock.Hex(),
		"l2BlockNum", l2BlockNum,
		"expiry", time.Unix(c.Expiry.Int64(), 0).Format(time.RFC1123Z),
		"statusEnum", c.Status,
	)
	log.Info("Attempting to defend pending L2 header challenge")

	tx, err := d.DefendL2Header(rblock, l2BlockNum)
	if err != nil {
		log.Error("Error defending L2 header challenge", "error", err)
		return fmt.Errorf("error defending L2 header challenge: %w", err)
	}

	log.Info("Pending L2 header challenge defended successfully", "tx", tx.Hash().Hex())

	return nil
}

// Defends an L2 header challenge by attempting to submit a header proof to the Challenge.sol contract.
func (d *Defender) DefendL2Header(rblock common.Hash, l2BlockNum *big.Int) (*types.Transaction, error) {
	// 1. Get the challenge key
	challengeHash, err := d.Ethereum.GetL2HeaderChallengeHash(rblock, l2BlockNum)
	if err != nil {
		return nil, fmt.Errorf("error getting challenge hash: %w", err)
	}

	// 2. Get the challenge
	challenge, err := d.Ethereum.GetL2HeaderChallenge(challengeHash)
	if err != nil {
		return nil, fmt.Errorf("error getting challenge: %w", err)
	}
	if challenge.Status != contracts.ChallengeL2HeaderStatusChallengerInitiated {
		return nil, fmt.Errorf("challenge is not pending")
	}

	// 3. Get the hashes of the header and previous header
	l2Block, err := d.LightLink.GetBlock(l2BlockNum.Uint64())
	if err != nil {
		return nil, fmt.Errorf("error getting block from l2: %w", err)
	}
	l2BlockHash := utils.HashWithoutExtraData(l2Block)

	l2PrevBlock, err := d.LightLink.GetBlock(l2BlockNum.Uint64() - 1)
	if err != nil {
		return nil, fmt.Errorf("error getting previous block from l2: %w", err)
	}
	l2PrevBlockHash := utils.HashWithoutExtraData(l2PrevBlock)

	// 4. Provide the headers
	tx, err := d.ProvideL2Header(challenge.Header.Rblock, l2BlockHash, false)
	if err != nil {
		return nil, fmt.Errorf("error providing header: %w", err)
	}
	d.Opts.Logger.Info("Provided header", "tx", tx.Hash().Hex(), "rblock", rblock.Hex(), "header", l2BlockHash.Hex())

	tx, err = d.ProvideL2Header(challenge.PrevHeader.Rblock, l2PrevBlockHash, false)
	if err != nil {
		return nil, fmt.Errorf("error providing previous header: %w", err)
	}
	d.Opts.Logger.Info("Provided previous header", "tx", tx.Hash().Hex(), "rblock", rblock.Hex(), "header", l2PrevBlockHash.Hex())

	// 5. Defend the challenge
	tx, err = d.Ethereum.DefendL2Header(challengeHash, l2BlockHash, l2PrevBlockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to defend l2 header challenge: %w", err)
	}

	return tx, nil
}

// Loads an L2 header from Celestia into the chainOracle.
func (d *Defender) ProvideL2Header(rblock common.Hash, l2Block common.Hash, skipShares bool) (*types.Transaction, error) {
	// check if the header is already provided
	headerProvided, _ := d.Ethereum.AlreadyProvidedHeader(l2Block)
	if headerProvided {
		d.Opts.Logger.Info("Header already provided", "block", rblock.Hex(), "header", l2Block.Hex())
		return types.NewTx(&types.LegacyTx{}), nil
	}

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundles, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	sharePointer, pointerIndex, err := node.FindHeaderSharesInBundles(bundles, l2Block, d.Namespace())
	if err != nil {
		return nil, fmt.Errorf("error finding header shares in the bundle: %w", err)
	}

	// Get proof the shares are in the bundle
	shareProof, err := d.Celestia.GetSharesProof(&node.CelestiaPointer{
		Height:     rheader.CelestiaPointers[pointerIndex].Height,
		ShareStart: rheader.CelestiaPointers[pointerIndex].ShareStart.Uint64(),
		ShareLen:   uint64(rheader.CelestiaPointers[pointerIndex].ShareLen),
	}, sharePointer)
	if err != nil {
		return nil, fmt.Errorf("error getting share proof: %w", err)
	}

	// Get proof the data is available
	celProof, err := d.GetDAProof(rblock, pointerIndex)
	if err != nil {
		return nil, fmt.Errorf("error proving data availability: %w", err)
	}

	// check if the shares are already provided
	provided, _ := d.Ethereum.AlreadyProvidedShares(rblock, shareProof.Data)

	attestationProof := chainOracleContract.AttestationProof{
		TupleRootNonce: celProof.RootNonce,
		Tuple: chainOracleContract.DataRootTuple{
			Height:   celProof.DataRootTuple.Height,
			DataRoot: celProof.DataRootTuple.DataRoot,
		},
		Proof: chainOracleContract.BinaryMerkleProof{
			SideNodes: celProof.Proof.SideNodes,
			Key:       celProof.Proof.Key,
			NumLeaves: celProof.Proof.NumLeaves,
		},
	}

	sp, err := contracts.NewShareProof(shareProof, attestationProof)
	if err != nil {
		return nil, fmt.Errorf("error creating share proof: %w", err)
	}

	// Provide the shares
	if !skipShares && !provided {
		tx, err := d.Ethereum.ProvideShares(rblock, pointerIndex, sp, attestationProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())

		// TODO: remove this sleep hack and fix Ethereum.Wait
		d.Opts.Logger.Info("Waiting for 10 seconds to ensure shares are available")
		time.Sleep(10 * time.Second)
		d.Ethereum.Wait(tx.Hash())
	}

	ranges := make([]chainOracleContract.ChainOracleShareRange, len(sharePointer.Ranges))
	for i, r := range sharePointer.Ranges {
		ranges[i] = chainOracleContract.ChainOracleShareRange{
			Start: big.NewInt(int64(r.Start)),
			End:   big.NewInt(int64(r.End)),
		}
	}

	// Finally, provide the header
	return d.Ethereum.ProvideHeader(rblock, shareProof.Data, ranges)
}

func (d *Defender) ProvideL2Tx(rblock common.Hash, l2Tx common.Hash, skipShares bool) (*types.Transaction, error) {
	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundles, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Find the transaction in the bundle
	sharePointer, pointerIndex, err := node.FindTxSharesInBundles(bundles, l2Tx, d.Namespace())
	if err != nil {
		return nil, fmt.Errorf("error finding tx shares in the bundle: %w", err)
	}

	// Get proof the shares are in the bundle
	shareProof, err := d.Celestia.GetSharesProof(&node.CelestiaPointer{
		Height:     rheader.CelestiaPointers[pointerIndex].Height,
		ShareStart: rheader.CelestiaPointers[pointerIndex].ShareStart.Uint64(),
		ShareLen:   uint64(rheader.CelestiaPointers[pointerIndex].ShareLen),
	}, sharePointer)
	if err != nil {
		return nil, fmt.Errorf("error getting share proof: %w", err)
	}

	// Get proof the data is available
	celProof, err := d.GetDAProof(rblock, pointerIndex)
	if err != nil {
		return nil, fmt.Errorf("error proving data availability: %w", err)
	}

	attestationProof := chainOracleContract.AttestationProof{
		TupleRootNonce: celProof.RootNonce,
		Tuple: chainOracleContract.DataRootTuple{
			Height:   celProof.DataRootTuple.Height,
			DataRoot: celProof.DataRootTuple.DataRoot,
		},
		Proof: chainOracleContract.BinaryMerkleProof{
			SideNodes: celProof.Proof.SideNodes,
			Key:       celProof.Proof.Key,
			NumLeaves: celProof.Proof.NumLeaves,
		},
	}

	sp, err := contracts.NewShareProof(shareProof, attestationProof)
	if err != nil {
		return nil, fmt.Errorf("error creating share proof: %w", err)
	}

	// Provide the shares
	if !skipShares {
		tx, err := d.Ethereum.ProvideShares(rblock, pointerIndex, sp, attestationProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())

		// TODO: remove this sleep hack and fix Ethereum.Wait
		d.Opts.Logger.Info("Waiting for 3 seconds to ensure shares are available")
		time.Sleep(10 * time.Second)
	}

	ranges := make([]chainOracleContract.ChainOracleShareRange, len(sharePointer.Ranges))
	for i, r := range sharePointer.Ranges {
		ranges[i] = chainOracleContract.ChainOracleShareRange{
			Start: big.NewInt(int64(r.Start)),
			End:   big.NewInt(int64(r.End)),
		}
	}

	// Finally, provide the transaction
	return d.Ethereum.ProvideLegacyTx(rblock, shareProof.Data, ranges)
}
