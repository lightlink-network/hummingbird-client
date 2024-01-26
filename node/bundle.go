package node

import (
	"hummingbird/utils"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxSizeLimit is the maximum size of a Celestia tx in bytes
const TxSizeLimit = uint64(1962441)

// Bundle is a collection of layer2 blocks which will be submitted to the
// data availability layer (Celestia).
type Bundle struct {
	Blocks []*types.Block
}

func (b *Bundle) Size() uint64 {
	return uint64(len(b.Blocks))
}

func (b *Bundle) Height() uint64 {
	return b.Blocks[len(b.Blocks)-1].Number().Uint64() + 1
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

// check if the bundle is under the celestia tx size
// limit of 1962441 bytes - 20000 bytes for tx overhead
func (b *Bundle) IsUnderTxLimit() (bool, uint64, uint64, error) {
	bundleEncoded, err := b.EncodeRLP()
	if err != nil {
		return false, 0, 0, err
	}
	bundleEncodedSize := uint64(len(bundleEncoded))
	bundleSizeLimit := TxSizeLimit - 20000
	if bundleEncodedSize > bundleSizeLimit {
		return false, bundleSizeLimit, bundleEncodedSize, nil
	}
	return true, bundleSizeLimit, bundleEncodedSize, nil
}

func (b *Bundle) Blob(namespace string) (*blob.Blob, error) {
	// 1. encode the bundle to RLP
	bundleRLP, err := b.EncodeRLP()
	if err != nil {
		return nil, err
	}

	// 2. get the blob
	return utils.BytesToBlob(namespace, bundleRLP)
}

type ShareRange struct {
	Start uint64
	End   uint64
}
