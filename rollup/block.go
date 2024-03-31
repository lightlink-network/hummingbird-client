package rollup

import (
	"hummingbird/node"

	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"

	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	*canonicalStateChainContract.CanonicalStateChainHeader
	Bundles []*node.Bundle
}

func (b *Block) CelestiaHeights() []uint64 {
	heights := make([]uint64, 0)
	for _, pointer := range b.CelestiaPointers {
		heights = append(heights, pointer.Height)
	}
	return heights
}

func (b *Block) L2Blocks() []*types.Block {
	blocks := []*types.Block{}
	for _, bundle := range b.Bundles {
		blocks = append(blocks, bundle.Blocks...)
	}
	return blocks
}
