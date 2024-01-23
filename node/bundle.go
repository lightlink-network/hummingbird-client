package node

import (
	"bytes"
	"fmt"
	"hummingbird/utils"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/celestia-node/blob"
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

// SharePointer is a pointer to some data inside a group of shares
type SharePointer struct {
	shares     []shares.Share
	StartShare int
	Ranges     []ShareRange
}

func NewSharePointer(_shares []shares.Share, startShare int, startIndex int, endShare int, endIndex int) *SharePointer {
	// add the start range
	ranges := []ShareRange{
		{
			Start: uint64(startIndex) + uint64(utils.ShareDataStart(_shares[startShare])),
			End:   uint64(_shares[startShare].Len()),
		},
	}

	// add the middle ranges
	for i := startShare + 1; i < endShare; i++ {
		r := ShareRange{
			Start: uint64(utils.ShareDataStart(_shares[i])),
			End:   uint64(_shares[i].Len()),
		}

		ranges = append(ranges, r)
	}

	// add the end range
	if startShare != endShare {
		ranges = append(ranges, ShareRange{
			Start: uint64(utils.ShareDataStart(_shares[endShare])),
		})
	}
	ranges[len(ranges)-1].End = uint64(endIndex) + uint64(utils.ShareDataStart(_shares[endShare]))

	return &SharePointer{
		shares:     _shares,
		StartShare: startShare,
		Ranges:     ranges,
	}
}

func (s *SharePointer) EndShare() int {
	return s.StartShare + len(s.Ranges) - 1
}

func (s *SharePointer) Bytes() []byte {
	data := []byte{}
	for i := 0; i < len(s.Ranges); i++ {
		data = append(data, s.shares[s.StartShare+i].ToBytes()[s.Ranges[i].Start:s.Ranges[i].End]...)
	}

	return data
}

func (s *SharePointer) Shares() []shares.Share {
	return s.shares[s.StartShare : s.EndShare()+1]
}

func (s *SharePointer) AllShares() []shares.Share {
	return s.shares
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

	return NewSharePointer(shares, startShare, startIndex, endShare, endIndex), nil
}
