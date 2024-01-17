package node

import (
	"bytes"
	"fmt"
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

type ShareRange struct {
	Start uint64
	End   uint64
}

type SharePointer struct {
	StartShare int
	StartIndex int
	EndShare   int
	EndIndex   int
}

// FinderHeaderShares finds the shares in the bundle which contain the header
func (b *Bundle) FindHeaderShares(hash common.Hash, namespace string) (*SharePointer, error) {
	// 1. find the block with the given hash
	var block *types.Block
	for _, b := range b.Blocks {
		if b.Hash() == hash {
			block = b
			break
		}
	}
	if block == nil {
		return nil, fmt.Errorf("block with hash %s not found in bundle", hash.Hex())
	}

	// 2. get the header RLP
	header := block.Header()
	headerRLP, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}

	// 3. Get the bundle shares
	bundleRLP, err := b.EncodeRLP()
	if err != nil {
		return nil, err
	}

	// 4. get bundle shares
	blob, err := utils.BytesToBlob(namespace, bundleRLP)
	if err != nil {
		return nil, err
	}
	shares, err := utils.BlobToShares(blob)
	if err != nil {
		return nil, err
	}

	// 5. find the header coords in the raw bundleRLP
	rlpStart := bytes.Index(bundleRLP, headerRLP)
	rlpEnd := rlpStart + len(headerRLP)

	if rlpStart == -1 {
		return nil, fmt.Errorf("encoded header not found in the bundle")
	}

	// 6. find the header coords in the shares
	startShare, startIndex := utils.RawIndexToSharesIndex(rlpStart, shares)
	endShare, endIndex := utils.RawIndexToSharesIndex(rlpEnd, shares)

	return &SharePointer{
		StartShare: startShare,
		StartIndex: startIndex,
		EndShare:   endShare,
		EndIndex:   endIndex,
	}, nil
}
