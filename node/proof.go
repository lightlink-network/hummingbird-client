package node

import (
	"context"

	"github.com/celestiaorg/celestia-app/pkg/square"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/tendermint/tendermint/types"
)

func (c *CelestiaClient) GetShareProofs(txHash []byte, pointer *SharePointer) (*types.ShareProof, error) {
	ctx := context.Background()

	// 1. Get the tx
	tx, err := c.trpc.Tx(ctx, txHash, true)
	if err != nil {
		return nil, err
	}

	// 2. Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &tx.Height)
	if err != nil {
		return nil, err
	}

	// 4. Get the share range inside the block
	shareRange, err := square.BlobShareRange(blockRes.Block.Data.Txs.ToSliceOfBytes(), int(tx.Index), 0, blockRes.Block.Header.Version.App)
	if err != nil {
		return nil, err
	}

	shareStart := uint64(shareRange.Start) //+ uint64(pointer.StartShare)
	shareEnd := uint64(shareRange.End)     //+ uint64(pointer.EndShare()+1)

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

func (c *CelestiaClient) GetShares(txHash []byte, namespace string) (share.NamespacedShares, error) {
	ctx := context.Background()

	ns, err := share.NewBlobNamespaceV0([]byte(c.Namespace()))
	if err != nil {
		return nil, err
	}

	// 1. Get the tx
	tx, err := c.trpc.Tx(ctx, txHash, true)
	if err != nil {
		return nil, err
	}

	// 2. Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &tx.Height)
	if err != nil {
		return nil, err
	}

	shareRange, err := square.TxShareRange(blockRes.Block.Data.Txs.ToSliceOfBytes(), int(tx.Index), blockRes.Block.Header.Version.App)
	if err != nil {
		return nil, err
	}

	h, err := c.client.Header.GetByHeight(ctx, uint64(tx.Height))
	if err != nil {
		return nil, err
	}

	shares, err := c.client.Share.GetSharesByNamespace(ctx, h, ns)
	if err != nil {
		return nil, err
	}

	return shares[shareRange.Start:shareRange.End], nil
}
