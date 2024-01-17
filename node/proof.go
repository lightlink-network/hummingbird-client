package node

import (
	"context"

	"github.com/celestiaorg/celestia-app/pkg/square"
	"github.com/tendermint/tendermint/types"
)

func (c *CelestiaClient) GetShareProofs(txHash []byte, pointer *SharePointer) (*types.ShareProof, error) {
	ctx := context.Background()

	// Get the tx
	tx, err := c.trpc.Tx(ctx, txHash, true)
	if err != nil {
		return nil, err
	}

	// Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &tx.Height)
	if err != nil {
		return nil, err
	}

	// Get the tx share range inside the block
	shareRange, err := square.TxShareRange(blockRes.Block.Data.Txs.ToSliceOfBytes(), int(tx.Index), blockRes.Block.Header.Version.App)
	if err != nil {
		return nil, err
	}

	shareStart := uint64(shareRange.Start) + uint64(pointer.StartShare)
	shareEnd := uint64(shareRange.Start) + uint64(pointer.EndShare)

	// Get the shares proof
	sharesProofs, err := c.trpc.ProveShares(ctx, uint64(tx.Height), shareStart, shareEnd)
	if err != nil {
		return nil, err
	}

	// Verify the shares proof
	if !sharesProofs.VerifyProof() {
		return nil, err
	}

	return &sharesProofs, nil
}
