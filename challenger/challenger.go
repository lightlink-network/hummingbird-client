package challenger

import (
	"hummingbird/node"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Opts struct {
	Logger *slog.Logger
	DryRun bool // DryRun indicates whether or not to actually submit the block to the L1 rollup contract.
}

type Challenger struct {
	*node.Node
	Opts *Opts
}

func NewChallenger(node *node.Node, opts *Opts) *Challenger {
	return &Challenger{Node: node, Opts: opts}
}

func (c *Challenger) ChallengeDA(index uint64, pointerIndex uint8) (*types.Transaction, common.Hash, error) {
	return c.Ethereum.ChallengeDataRootInclusion(index, pointerIndex)
}
