package rollup

import (
	"fmt"
	"hummingbird/node"
	"hummingbird/node/contracts"
	"hummingbird/utils"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type Opts struct {
	BundleSize uint64        // BundleSize is the number of blocks to include in each bundle.
	PollDelay  time.Duration // PollDelay is the time to wait between polling for new blocks.
	Logger     *slog.Logger
	DryRun     bool // DryRun indicates whether or not to actually submit the block to the L1 rollup contract.

	StoreCelestiaPointers bool // StoreCelestiaPointers indicates whether or not to store the Celestia pointers in the local database.
	StoreHeaders          bool // StoreHeaders indicates whether or not to store the rollup headers in the local database.
}

type Rollup struct {
	*node.Node
	Opts *Opts
}

func NewRollup(n *node.Node, opts *Opts) *Rollup {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	log := opts.Logger.With("func", "NewRollup")

	// Check if local storage should be disabled
	// if so warn the user and disable it.
	disableStorage := false
	if n.Store == nil {
		log.Warn("A Store was not set on the Node, disabling local storage")
		disableStorage = true
	}
	if opts.DryRun {
		log.Warn("DryRun is enabled, disabling local storage")
		disableStorage = true
	}
	if disableStorage {
		opts.StoreCelestiaPointers = false
		opts.StoreHeaders = false
	}

	return &Rollup{Node: n, Opts: opts}
}

// CreateNextBlock creates a new rollup block from the current state of the
// lightlink network and pushes it to the data availability layer (Celestia).
// It returns the new block and an error if one occurred.
//
// Note: This function does not submit the block to the L1 rollup contract.
// See: Rollup.SubmitBlock
func (r *Rollup) CreateNextBlock() (*Block, error) {
	// 0. fetch the current epoch = eth height
	epoch, err := r.Ethereum.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to get current epoch: %w", err)
	}

	// 1. fetch ll height
	llHeight, err := r.LightLink.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to get current llheight: %w", err)
	}

	// 2. fetch the last rollup header
	head, err := r.Ethereum.GetRollupHead()
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to get current head: %w", err)
	}

	// 3. calculate bundle size
	bundleSize := r.Opts.BundleSize
	if llHeight-head.L2Height < bundleSize {
		bundleSize = llHeight - head.L2Height
	}

	// 4. calc prevHash from the last rollup header
	prevHash, err := contracts.HashCanonicalStateChainHeader(&head)
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to get current prevHash: %w", err)
	}

	// 5. fetch the next bundle of blocks from ll
	l2blocks, err := r.LightLink.GetBlocks(head.L2Height+1, head.L2Height+1+bundleSize)
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to get l2blocks: %w", err)
	}
	bundle := &node.Bundle{l2blocks}

	// 6. upload the bundle to celestia
	pointer, err := r.Celestia.PublishBundle(*bundle)
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to publish bundle: %w", err)
	}

	// 7. create the rollup header
	header := &contracts.CanonicalStateChainHeader{
		Epoch:            epoch,
		L2Height:         head.L2Height + bundleSize,
		PrevHash:         prevHash,
		TxRoot:           bundle.TxRoot(),
		BlockRoot:        bundle.BlockRoot(),
		StateRoot:        bundle.StateRoot(),
		CelestiaHeight:   pointer.Height,
		CelestiaDataRoot: pointer.DataRoot,
	}

	// 8. calculate the hash of the header
	hash, err := contracts.HashCanonicalStateChainHeader(header)
	if err != nil {
		return nil, fmt.Errorf("createNextBlock: Failed to hash header: %w", err)
	}

	// 9. Optionally store the header in the local database
	if r.Opts.StoreHeaders {
		if err := r.Node.Store.Put(hash[:], utils.MustJsonMarshal(header)); err != nil {
			return nil, err
		}
	}

	// 10. Optionally store the Celestia pointer in the local database
	// Required for the Celestia proof.
	if r.Opts.StoreCelestiaPointers {
		if err := r.Node.Store.Put(pointer.TxHash[:], utils.MustJsonMarshal(pointer)); err != nil {
			return nil, err
		}
	}

	return &Block{header, bundle, pointer}, nil
}

func (b *Rollup) SubmitBlock(block *Block) (*types.Transaction, error) {
	log := b.Opts.Logger.With("func", "SubmitBlock")

	tx, err := b.Ethereum.PushRollupHead(block.CanonicalStateChainHeader)
	if err != nil {
		log.Error("Failed to push rollup head", "error", err)
		return nil, err
	}

	log.Info("Submitted rollup block", "tx", tx.Hash().Hex(), "epoch", block.Epoch, "l2Height", block.L2Height, "celestiaHeight", block.CelestiaHeight)
	return tx, nil
}

func (r *Rollup) CreateAndSubmitNextBlock() (*Block, uint64, error) {
	log := r.Opts.Logger.With("func", "CreateAndSubmitNextBlock")

	// 1. create the next rollup block
	block, err := r.CreateNextBlock()
	if err != nil {
		log.Error("Failed to create next block", "error", err)
		return nil, 0, err
	}

	// 2. submit the block to the rollup contract
	tx, err := r.SubmitBlock(block)
	if err != nil {
		log.Error("Failed to submit block", "error", err)
		return nil, 0, err
	}

	// 3. wait for the tx
	receipt, err := r.Ethereum.Wait(tx.Hash())
	if err != nil {
		log.Error("Failed to wait for tx", "error", err)
		return nil, 0, err
	}

	if receipt.Status != 1 {
		log.Error("Transaction failed", "tx", tx.Hash().Hex(), "status", receipt.Status)
		return nil, 0, err
	}

	// 3. wait for the block to be mined
	h, err := r.Ethereum.GetRollupHeight()
	if err != nil {
		log.Error("Failed to get rollup height", "error", err)
		return nil, 0, err
	}

	log.Info("Rollup chain updated", "rollup_l2height", block.L2Height, "bundle_size", len(block.Blocks), "rollup_height", h, "epoch", block.Epoch, "tx", receipt.TxHash.Hex(), "gas_used", receipt.GasUsed)
	return block, h, nil
}

func (r *Rollup) Run() error {
	log := r.Opts.Logger.With("func", "Run")

	// get last rollup height
	head, err := r.Ethereum.GetRollupHead()
	if err != nil {
		log.Error("Failed to get rollup height", "error", err)
		return err
	}
	log.Info("Starting rollup", "rollup_ll_height", head.L2Height, "rollup_ll_epoch", head.Epoch)

	for {
		// 1. get next rollup target height
		target, err := r.nextRollupTarget()
		if err != nil {
			log.Error("Failed to get next rollup target", "error", err)
			return err
		}

		// 2. wait for the target height to be reached
		err = r.awaitL2Height(target)
		if err != nil {
			log.Error("Failed to await L2 height", "error", err)
			return err
		}

		// 3. create the next rollup block
		block, err := r.CreateNextBlock()
		if err != nil {
			log.Error("Failed to create next block", "error", err)
			return err
		}

		// 4. submit the block to the rollup contract
		tx, err := r.SubmitBlock(block)
		if err != nil {
			log.Error("Failed to submit block", "error", err)
			return err
		}

		hash, _ := contracts.HashCanonicalStateChainHeader(block.CanonicalStateChainHeader)
		log.Info("Submitted rollup block",
			"tx", tx.Hash().Hex(),
			"hash", hash,
			"epoch", block.Epoch,
			"l2Height", block.L2Height,
			"celestiaHeight", block.CelestiaHeight,
			"daTx", block.CelestiaPointer.TxHash.Hex(),
		)

	}

}

// returns the layer2 block height that will trigger the next rollup
func (r *Rollup) nextRollupTarget() (uint64, error) {
	log := r.Opts.Logger.With("func", "nextRollupTarget")

	head, err := r.Ethereum.GetRollupHead()
	if err != nil {
		log.Error("Failed to get rollup height", "error", err)
		return 0, err
	}

	// get the next rollup target
	return head.L2Height + r.Opts.BundleSize, nil
}

func (r *Rollup) awaitL2Height(h uint64) error {

	for {
		llHeight, err := r.LightLink.GetHeight()
		if err != nil {
			return fmt.Errorf("failed to get layer 2 height: %w", err)
		}

		if llHeight > h {
			return nil
		}

		time.Sleep(r.Opts.PollDelay)
	}
}
