package defender

import (
	"fmt"
	"hummingbird/node"
	"time"

	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	challengeContract "hummingbird/node/contracts/Challenge.sol"
)

type Opts struct {
	Logger      *slog.Logger
	WorkerDelay time.Duration
	DryRun      bool // DryRun indicates whether or not to actually submit the block to the L1 rollup contract.
}

type Defender struct {
	*node.Node
	Opts *Opts
}

func NewDefender(node *node.Node, opts *Opts) *Defender {
	return &Defender{Node: node, Opts: opts}
}

func (d *Defender) Start() error {
	d.retryMissedDAChallenges()
	go d.retryActiveDAChallengesWorker()

	if err := d.WatchAndDefendDAChallenges(); err != nil {
		return fmt.Errorf("error watching and defending DA challenges: %w", err)
	}
	return nil
}

func (d *Defender) WatchAndDefendDAChallenges() error {
	challenges := make(chan *challengeContract.ChallengeChallengeDAUpdate)
	subscription, err := d.Ethereum.WatchChallengesDA(challenges)
	if err != nil {
		return fmt.Errorf("error starting WatchChallengesDA: %w", err)
	}
	defer subscription.Unsubscribe()

	d.Opts.Logger.Info("Defender is watching for DA challenges")

	// Listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		close(challenges)
	}()

	var wg sync.WaitGroup

	for challenge := range challenges {
		wg.Add(1)
		go func(challenge *challengeContract.ChallengeChallengeDAUpdate) {
			defer wg.Done()

			header, err := d.Ethereum.GetRollupHeaderByHash(challenge.BlockHash)
			if err != nil {
				d.Opts.Logger.Error("error getting rollup header by hash:", "error", err)
			}
			err = d.Store.StoreLastScannedBlockNumber(header.Epoch)
			if err != nil {
				d.Opts.Logger.Error("error storing last scanned block number:", "error", err)
			}

			err = d.handleDAChallenge(challenge)
			if err != nil {
				d.Opts.Logger.Error("error handling challenge:", "challenge", challenge, "error", err)
				err := d.Store.StoreActiveDAChallenge(challenge)
				if err != nil {
					d.Opts.Logger.Error("error storing active DA challenge:", "challenge", challenge, "error", err)
				}
			}
		}(challenge)
	}

	// Wait for all challenges to be handled before returning
	wg.Wait()

	return nil
}

func (d *Defender) handleDAChallenge(challenge *challengeContract.ChallengeChallengeDAUpdate) error {
	blockHash := common.BytesToHash(challenge.BlockHash[:])

	d.Opts.Logger.Info("DA challenge received", "block", blockHash.Hex(), "block_index", challenge.BlockIndex, "expiry", challenge.Expiry, "status", challenge.Status)

	if challenge.Status != 1 {
		return nil
	}

	celestiaTx, err := d.GetDAPointer(challenge.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting CelestiaTx: %w", err)
	}

	if celestiaTx == nil {
		d.Opts.Logger.Info("No CelestiaTx found", "block:", blockHash.Hex())
		return nil
	}

	d.Opts.Logger.Info("Found CelestiaTx", "tx_hash", celestiaTx.TxHash.Hex(), "block_hash", blockHash.Hex())

	tx, err := d.DefendDA(challenge.BlockHash)
	if err != nil {
		return fmt.Errorf("error defending DA: %w", err)
	}

	d.Opts.Logger.Info("DA challenge defended", "tx", tx.Hash().Hex(), "block", blockHash.Hex(), "block_index", challenge.BlockIndex, "expiry", challenge.Expiry, "status", challenge.Status)
	return nil
}

func (d *Defender) DefendDA(block common.Hash) (*types.Transaction, error) {
	proof, err := d.ProveDA(block)
	if err != nil {
		return nil, fmt.Errorf("failed to prove data availability: %w", err)
	}

	d.Opts.Logger.Debug("Submitting data availability proof to L1 rollup contract", "block", block.Hex(), "dataroot", hexutil.Encode(proof.Tuple.DataRoot[:]))
	return d.Ethereum.DefendDataRootInclusion(block, proof)
}

func (d *Defender) ProveDA(block common.Hash) (*node.CelestiaProof, error) {
	pointer, err := d.GetDAPointer(block)
	if err != nil {
		return nil, fmt.Errorf("failed to get Celestia pointer: %w", err)
	}

	if pointer == nil {
		return nil, fmt.Errorf("no Celestia pointer found")
	}

	return d.Celestia.GetProof(pointer)
}

func (d *Defender) retryActiveDAChallengesWorker() {
	ticker := time.NewTicker(d.Opts.WorkerDelay)
	defer ticker.Stop()

	for range ticker.C {
		d.Opts.Logger.Info("Retrying active DA challenges...")
		challenges, err := d.Store.GetActiveDAChallenges()
		if err != nil {
			d.Opts.Logger.Error("error getting active DA challenges from store:", "error", err)
			continue
		}
		for _, challenge := range challenges {
			block := common.BytesToHash(challenge.BlockHash[:])

			// Check if challenge has expired, if so delete from active challenges and continue
			if challenge.Expiry.Int64() <= time.Now().Unix() {
				d.Opts.Logger.Info("Active DA challenge has expired, deleting from active challenges", "challengeBlock", block, "expiry", challenge.Expiry)
				err = d.Store.DeleteActiveDAChallenge(challenge.BlockHash)
				if err != nil {
					d.Opts.Logger.Error("error deleting active DA challenge:", "challengeBlock", block, "error", err)
				}
				continue
			}

			err = d.handleDAChallenge(challenge)
			if err != nil {
				d.Opts.Logger.Error("error retrying active DA challenge:", "challengeBlock", block, "error", err)
				continue
			}

			err = d.Store.DeleteActiveDAChallenge(challenge.BlockHash)
			if err != nil {
				d.Opts.Logger.Error("error deleting active DA challenge:", "challengeBlock", block, "error", err)
			}
		}
		d.Opts.Logger.Info("Active DA challenges retry worker finished")
	}
}

func (d *Defender) retryMissedDAChallenges() {
	d.Opts.Logger.Info("Retrying missed DA challenges...")
	lastScannedBlockNumber, err := d.Store.GetLastScannedBlockNumber()
	if err != nil {
		d.Opts.Logger.Error("error getting last scanned block number:", "error", err)
		return
	}

	opts := &bind.FilterOpts{
		Start: lastScannedBlockNumber,
	}

	status := []uint8{1}

	challenges, err := d.Ethereum.FilterChallengeDAUpdate(opts, nil, nil, status)
	if err != nil {
		d.Opts.Logger.Error("error filtering challenges:", "error", err)
		return
	}

	// iterate through challenges and retry
	for challenges.Next() {
		challenge := challenges.Event // Access the current challenge

		// Check if challenge has already been defended
		challengeInfo, err := d.Ethereum.GetDataRootInclusionChallenge(challenge.BlockHash)
		if err != nil {
			d.Opts.Logger.Error("error getting data root inclusion challenge:", "error", err)
			continue
		}

		if challengeInfo.Status != 1 {
			d.Opts.Logger.Info("DA challenge has already been defended", "blockIndex", challenge.BlockIndex)
			continue
		}

		err = d.handleDAChallenge(challenge)
		if err != nil {
			d.Opts.Logger.Error("error retrying missed DA challenge:", "blockIndex", challenge.BlockIndex, "error", err)
			continue
		}
	}
	d.Opts.Logger.Info("Missed DA challenges retry finished")
}

func (d *Defender) ProvideL2Header(rblock common.Hash, l2Block common.Hash) (*types.Transaction, error) {

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundle, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Get a pointer to the shares that contain the header
	sharePointer, err := bundle.FindHeaderShares(l2Block, d.Namespace())
	if err != nil {
		return nil, fmt.Errorf("error finding header shares in the bundle: %w", err)
	}

	// Get proof the shares are in the bundle
	shareProof, err := d.Celestia.GetSharesProof(&node.CelestiaPointer{
		Height:     rheader.CelestiaHeight,
		ShareStart: rheader.CelestiaShareStart,
		ShareLen:   rheader.CelestiaShareLen,
	}, sharePointer)
	if err != nil {
		return nil, fmt.Errorf("error getting share proof: %w", err)
	}

	// Get proof the data is available
	celProof, err := d.ProveDA(rblock)
	if err != nil {
		return nil, fmt.Errorf("error proving data availability: %w", err)
	}

	// Provide the shares
	tx, err := d.Ethereum.ProvideShares(rblock, shareProof, celProof)
	if err != nil {
		return nil, fmt.Errorf("error providing shares: %w", err)
	}
	d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
	d.Ethereum.Wait(tx.Hash())

	// Finally, provide the header
	return d.Ethereum.ProvideHeader(rblock, shareProof.Data, *sharePointer)
}
