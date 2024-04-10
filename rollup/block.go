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

func (b *Block) GetCelestiaPointers() []*node.CelestiaPointer {
	ps := make([]*node.CelestiaPointer, 0)
	for _, pointer := range b.CelestiaPointers {
		ps = append(ps, &node.CelestiaPointer{
			Height:     pointer.Height,
			ShareStart: pointer.ShareStart.Uint64(),
			ShareLen:   uint64(pointer.ShareLen),
		})
	}
	return ps
}
