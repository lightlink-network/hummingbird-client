package defender

import (
	"fmt"
	"hummingbird/node"
	"strings"
	"time"

	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"hummingbird/node/contracts"
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
//
// Runs every d.Opts.WorkerDelay ms to scan Challenge.sol event log, find any pending
// DA challenges (status 1) and attempt to defend them.
func (d *Defender) Start() error {
	// Run once at initialization
	if err := d.findAndDefendChallenges(); err != nil {
		d.Opts.Logger.Debug("error defending historic DA challenges: %w", err)
	}

	ticker := time.NewTicker(d.Opts.WorkerDelay)
	defer ticker.Stop()

	// Run every d.Opts.WorkerDelay ms
	for range ticker.C {
		if err := d.findAndDefendChallenges(); err != nil {
			d.Opts.Logger.Debug("error defending historic DA challenges: %w", err)
			continue
		}
	}
	return nil
}

// Scans the Challenge.sol contract for pending DA challenges and attempts to defend them.
//
// Begins scanning from the first rollup headers epoch and iterates through all historic
// challenges, attempting to defend each one if it is still pending (status 1).
func (d *Defender) findAndDefendChallenges() error {
	d.Opts.Logger.Debug("Starting log scan for historic pending DA challenges")

	challenges, err := d.Ethereum.FilterChallengeDAUpdate(nil, nil, nil, []uint8{contracts.ChallengeDAStatusChallengerInitiated})
	if err != nil {
		return fmt.Errorf("error filtering challenges: %w", err)
	}

	// iterate through historic challenges events
	for challenges.Next() {
		challenge := challenges.Event
		blockHash := common.BytesToHash(challenge.BlockHash[:])

		// check if challenge has already been defended by getting the current status
		// required as we are scanning historic logs, and the challenge may have been defended since log was emitted
		challengeInfo, err := d.Ethereum.GetDataRootInclusionChallenge(blockHash)
		if err != nil {
			d.Opts.Logger.Error("error getting data root inclusion challenge", "block_hash", blockHash, "error", err)
			continue
		}

		// we are only interested in challenges that have been initiated by a challenger, ready to be defended
		if challengeInfo.Status != contracts.ChallengeDAStatusChallengerInitiated {
			continue
		}

		err = d.handleDAChallenge(challenge)
		if err != nil {
			d.Opts.Logger.Error("error handling DA challenge", "block_hash", blockHash, "error", err)
			continue
		}

	}
	d.Opts.Logger.Debug("Finished log scan for historic pending DA challenges")

	return nil
}

// Handles a DA challenge by attempting to defend it.
//
// If the challenged data root is not yet available, it will be ignored
// and retried later by the findAndDefendChallenges() worker function.
func (d *Defender) handleDAChallenge(challenge *challengeContract.ChallengeChallengeDAUpdate) error {
	blockHash := common.BytesToHash(challenge.BlockHash[:])
	statusString := contracts.DAChallengeStatusToString(challenge.Status)

	log := d.Opts.Logger.With(
		"blockHash", blockHash.Hex(),
		"blockIndex", challenge.BlockIndex,
		"expiry", time.Unix(challenge.Expiry.Int64(), 0).Format(time.RFC1123Z),
		"statusEnum", challenge.Status,
		"statusString", statusString,
	)
	log.Info("Pending DA challenge log event found")

	// attempt to defend the challenge by submitting a tx to the Challenge contract
	tx, err := d.DefendDA(challenge.BlockHash, uint8(challenge.PointerIndex.Uint64()))
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
	return d.Ethereum.DefendDataRootInclusion(block, proof)
}

// Gets the Celestia pointer for the given block hash and queries Celestia for a proof
// of data availability.
func (d *Defender) GetDAProof(block common.Hash, pointerIndex uint8) (*node.CelestiaProof, error) {
	pointers, err := d.GetDAPointer(block)
	if err != nil {
		return nil, fmt.Errorf("failed to get Celestia pointer: %w", err)
	}
	if pointers == nil {
		return nil, fmt.Errorf("no Celestia pointer found")
	}
	return d.Celestia.GetProof(pointers[pointerIndex])
}

func (d *Defender) ProvideL2Header(rblock common.Hash, pointerIndex uint8, l2Block common.Hash, skipShares bool) (*types.Transaction, error) {

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundles, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Get the bundle that contains the header
	bundle := bundles[pointerIndex]

	// Get a pointer to the shares that contain the header
	sharePointer, err := bundle.FindHeaderShares(l2Block, d.Namespace())
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

	// Provide the shares
	if !skipShares {
		tx, err := d.Ethereum.ProvideShares(rblock, pointerIndex, shareProof, celProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())

		// TODO: remove this sleep hack and fix Ethereum.Wait
		d.Opts.Logger.Info("Waiting for 3 seconds to ensure shares are available")
		time.Sleep(10 * time.Second)
	}

	// Finally, provide the header
	return d.Ethereum.ProvideHeader(rblock, shareProof.Data, *sharePointer)
}

func (d *Defender) ProvideL2Tx(rblock common.Hash, pointerIndex uint8, l2Tx common.Hash, skipShares bool) (*types.Transaction, error) {

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundles, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Get the bundle that contains the transaction
	bundle := bundles[pointerIndex]

	// Get a pointer to the shares that contain the transaction
	sharePointer, err := bundle.FindTxShares(l2Tx, d.Namespace())
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

	// Provide the shares
	if !skipShares {
		tx, err := d.Ethereum.ProvideShares(rblock, pointerIndex, shareProof, celProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())

		// TODO: remove this sleep hack and fix Ethereum.Wait
		d.Opts.Logger.Info("Waiting for 3 seconds to ensure shares are available")
		time.Sleep(10 * time.Second)
	}

	// Finally, provide the transaction
	return d.Ethereum.ProvideLegacyTx(rblock, shareProof.Data, *sharePointer)
}
