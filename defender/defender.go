package defender

import (
	"fmt"
	"hummingbird/node"

	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	challengeContract "hummingbird/node/contracts/Challenge.sol"
)

type Opts struct {
	Logger *slog.Logger
	DryRun bool // DryRun indicates whether or not to actually submit the block to the L1 rollup contract.
}

type Defender struct {
	*node.Node
	Opts *Opts
}

func NewDefender(node *node.Node, opts *Opts) *Defender {
	return &Defender{Node: node, Opts: opts}
}

func (d *Defender) Start() error {
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
			err := d.handleDAChallenge(challenge)
			if err != nil {
				d.Opts.Logger.Error("error handling challenge:", "challenge", challenge, "error", err)
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

	celestiaTx, err := d.Store.GetDAPointer(challenge.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting CelestiaTx: %w", err)
	}

	if celestiaTx == nil {
		d.Opts.Logger.Info("No CelestiaTx found", "block:", blockHash.Hex())
		return nil
	}

	d.Opts.Logger.Info("Found CelestiaTx", "tx_hash", celestiaTx.TxHash.Hex(), "block_hash", blockHash.Hex())

	tx, err := d.DefendDA(challenge.BlockHash, celestiaTx.TxHash)
	if err != nil {
		return fmt.Errorf("error defending DA: %w", err)
	}

	d.Opts.Logger.Info("DA challenge defended", "tx", tx.Hash().Hex(), "block", blockHash.Hex(), "block_index", challenge.BlockIndex, "expiry", challenge.Expiry, "status", challenge.Status)
	return nil
}

func (d *Defender) DefendDA(block common.Hash, txHash common.Hash) (*types.Transaction, error) {
	proof, err := d.ProveDA(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to prove data availability: %w", err)
	}

	d.Opts.Logger.Debug("Submitting data availability proof to L1 rollup contract", "block", block.Hex(), "dataroot", hexutil.Encode(proof.Tuple.DataRoot[:]))
	return d.Ethereum.DefendDataRootInclusion(block, proof)
}

func (d *Defender) ProveDA(txHash common.Hash) (*node.CelestiaProof, error) {
	return d.Celestia.GetProof(txHash[:])
}
