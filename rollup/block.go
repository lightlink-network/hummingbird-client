package rollup

import (
	"hummingbird/node"

	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
)

type Block struct {
	*canonicalStateChainContract.CanonicalStateChainHeader
	*node.Bundle
}

func (b *Block) CelestiaHeights() []uint64 {
	heights := make([]uint64, 0)
	for _, pointer := range b.CelestiaPointers {
		heights = append(heights, pointer.Height)
	}
	return heights
}
