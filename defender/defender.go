package defender

import (
	"fmt"
	"hummingbird/node"
	"hummingbird/node/contracts"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

func (d *Defender) ProveDA(txHash common.Hash) (*node.CelestiaProof, error) {
	return d.Celestia.GetProof(txHash[:])
}

func (d *Defender) DefendDA(block common.Hash, txHash common.Hash) (*types.Transaction, error) {
	proof, err := d.ProveDA(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to prove data availability: %w", err)
	}

	d.Opts.Logger.Debug("Submitting data availability proof to L1 rollup contract", "block", block.Hex(), "dataroot", hexutil.Encode(proof.Tuple.DataRoot[:]))
	return d.Ethereum.DefendDataRootInclusion(block, proof)
}

func (d *Defender) Start() error {
	challenges := make(chan *contracts.ChallengeContractChallengeDAUpdate)
	subscription, err := d.Ethereum.WatchChallengesDA(challenges)
	if err != nil {
		return fmt.Errorf("error starting WatchChallengesDA: %w", err)
	}
	defer subscription.Unsubscribe()

	d.Opts.Logger.Info("Defender started and watching for challenges")

	for challenge := range challenges {
		if challenge.Status != 1 {
			continue
		}

		blockHash := common.BytesToHash(challenge.BlockHash[:])

		d.Opts.Logger.Info("DA challenge received", "block", blockHash.Hex(), "block_index", challenge.BlockIndex, "expiry", challenge.Expiry, "status", challenge.Status)

		celestiaTx, err := d.Store.GetDAPointer(challenge.BlockHash)
		if err != nil {
			d.Opts.Logger.Error("error getting CelestiaTx:", err)
			continue
		}

		if celestiaTx == nil {
			d.Opts.Logger.Info("no CelestiaTx found", "block:", blockHash.Hex())
			continue
		}

		d.Opts.Logger.Info("Found CelestiaTx", "tx_hash", celestiaTx.TxHash.Hex(), "block_hash", blockHash.Hex())

		tx, err := d.DefendDA(challenge.BlockHash, celestiaTx.TxHash)
		if err != nil {
			d.Opts.Logger.Error("error defending DA:", err)
		}

		d.Opts.Logger.Info("DA challenge defended", "tx", tx.Hash().Hex(), "block", blockHash.Hex(), "block_index", challenge.BlockIndex, "expiry", challenge.Expiry, "status", challenge.Status)
	}

	return nil
}
