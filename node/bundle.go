package node

import (
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Bundle is a collection of layer2 blocks which will be submitted to the
// data availability layer (Celestia).
type Bundle struct {
	Blocks []*types.Block
}

func (b *Bundle) EncodeRLP() ([]byte, error) {
	return rlp.EncodeToBytes(&b.Blocks)
}

func (b *Bundle) DecodeRLP(data []byte) error {
	return rlp.DecodeBytes(data, &b.Blocks)
}

// get the root of the merkle tree containing all the blocks in the bundle
func (b *Bundle) BlockRoot() common.Hash {
	hashes := make([]common.Hash, len(b.Blocks))
	for i, block := range b.Blocks {
		hashes[i] = block.Hash()
	}

	return utils.CalculateMerkleRoot(hashes...)
}

// get the root of the merkle tree containing all the transactions, in all the blocks, in the bundle
func (b *Bundle) TxRoot() common.Hash {
	hashes := []common.Hash{}

	for _, block := range b.Blocks {
		for _, tx := range block.Transactions() {
			hashes = append(hashes, tx.Hash())
		}
	}

	return utils.CalculateMerkleRoot(hashes...)
}

// get the stateroot of the last block in the bundle
func (b *Bundle) StateRoot() common.Hash {
	last := b.Blocks[len(b.Blocks)-1]
	return last.Header().Root
}
