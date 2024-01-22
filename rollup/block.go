package rollup

import (
	"hummingbird/node"

	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
)

type Block struct {
	*canonicalStateChainContract.CanonicalStateChainHeader
	*node.Bundle
	*node.CelestiaPointer
}
