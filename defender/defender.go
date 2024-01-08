package defender

import (
	"encoding/json"
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

func (d *Defender) ProveDAByBlockHash(blockHash common.Hash) (*node.CelestiaProof, error) {
	key := append([]byte("da_pointer"), blockHash[:]...)
	buf, err := d.Node.Store.Get(key)
	if err != nil {
		return nil, err
	}

	var pointer node.CelestiaPointer
	err = json.Unmarshal(buf, &pointer)
	if err != nil {
		return nil, err
	}

	return d.ProveDA(pointer.TxHash)
}
