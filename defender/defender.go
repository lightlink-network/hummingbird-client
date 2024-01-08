package defender

import (
	"hummingbird/node"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
)

type Opts struct {
	Logger *slog.Logger
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
