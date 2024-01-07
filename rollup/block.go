package rollup

import (
	"hummingbird/node"
	"hummingbird/node/contracts"
)

type Block struct {
	*contracts.CanonicalStateChainHeader
	*node.Bundle
	*node.CelestiaPointer
}
