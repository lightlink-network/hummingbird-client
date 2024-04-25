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
	"github.com/tendermint/tendermint/crypto/merkle"
)

// TxSizeLimit is the maximum size of a Celestia tx in bytes
const TxSizeLimit = uint64(1962441)

// Bundle is a collection of layer2 blocks which will be submitted to the
// data availability layer (Celestia).
type Bundle struct {
	Blocks []*types.Block
}

func NewBundleFromShares(s []shares.Share) (*Bundle, error) {
	// 1. extract the raw data from the shares
	data := []byte{}
	for _, share := range s {
		d, err := share.RawData()
		if err != nil {
			return nil, err
		}

		data = append(data, d...)
	}
	// 2. decode the bundle from RLP
	bundle := &Bundle{}
	return bundle, bundle.DecodeRLP(data)
}

func (b *Bundle) Size() uint64 {
	return uint64(len(b.Blocks))
}

func (b *Bundle) Height() uint64 {
	return b.Blocks[len(b.Blocks)-1].Number().Uint64()
}

func (b *Bundle) EncodeRLP() ([]byte, error) {
	return rlp.EncodeToBytes(&b.Blocks)
}

func (b *Bundle) DecodeRLP(data []byte) error {
	size := utils.RlpNextItemSize(data)
	if size == -1 {
		return fmt.Errorf("DecodeRLP: invalid rlp data")
	}

	return rlp.DecodeBytes(data[:size], &b.Blocks)
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
// limit of 1962441 bytes - 196245 bytes (10%) for tx overhead
func (b *Bundle) IsUnderTxLimit() (bool, uint64, uint64, error) {
	bundleEncoded, err := b.EncodeRLP()
	if err != nil {
		return false, 0, 0, err
	}
	bundleEncodedSize := uint64(len(bundleEncoded))
	bundleSizeLimit := TxSizeLimit - 196245
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

func (b *Bundle) Shares(namespace string) ([]shares.Share, error) {
	// 1. get the blob
	blob, err := b.Blob(namespace)
	if err != nil {
		return nil, err
	}

	// 2. get the shares
	return utils.BlobToShares(blob)
}

// FinderHeaderShares finds the shares in the bundle which contain the header
// Returns a pointer to the data in the shares
func (b *Bundle) FindHeaderShares(hash common.Hash, namespace string) (*SharePointer, error) {
	// 1. find the block with the given hash
	var block *types.Block
	for _, b := range b.Blocks {
		// This is a hack to fix the issue with the extra data in the block header.
		// TODO: remove this hack and fix extra data before bundle upload
		if utils.HashWithoutExtraData(b).Hex() == hash.Hex() {
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

	// 3. Get the bundle rlp
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

// FinderTxShares finds the shares in the bundle which contain the transaction
// Returns a pointer to the data in the shares
func (b *Bundle) FindTxShares(hash common.Hash, namespace string) (*SharePointer, error) {
	// 1. find the tx with the given hash
	var tx *types.Transaction
	for _, block := range b.Blocks {
		for _, t := range block.Transactions() {
			if t.Hash().Hex() == hash.Hex() {
				tx = t
				break
			}
		}
	}

	if tx == nil {
		return nil, fmt.Errorf("tx with hash %s not found in bundle", hash.Hex())
	}

	if tx.Type() != types.LegacyTxType {
		return nil, fmt.Errorf("tx with hash %s is not a legacy tx, only legacy txns supported", hash.Hex())
	}

	// 2. get the tx RLP
	txRLP, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	// 3. Get the bundle rlp
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

	// 4. find the tx coords in the raw bundleRLP
	rlpStart := bytes.Index(bundleRLP, txRLP)
	rlpEnd := rlpStart + len(txRLP)

	if rlpStart == -1 {
		return nil, fmt.Errorf("encoded tx not found in the bundle")
	}

	// 5. find the tx coords in the shares
	startShare, startIndex := utils.RawIndexToSharesIndex(rlpStart, shares)
	endShare, endIndex := utils.RawIndexToSharesIndex(rlpEnd, shares)

	return NewSharePointer(shares, startShare, startIndex, endShare, endIndex), nil

	// TODO: this code repeats the same logic as FindHeaderShares, we should refactor it
	// to avoid code duplication. `FindBytesShares` ?
}

func FindHeaderSharesInBundles(bundles []*Bundle, hash common.Hash, namespace string) (*SharePointer, uint8, error) {
	for i, bundle := range bundles {
		pointer, err := bundle.FindHeaderShares(hash, namespace)
		if err == nil {
			return pointer, uint8(i), nil
		}
	}
	return nil, 0, fmt.Errorf("header not found in any bundle")
}

func FindTxSharesInBundles(bundles []*Bundle, hash common.Hash, namespace string) (*SharePointer, uint8, error) {
	for i, bundle := range bundles {
		pointer, err := bundle.FindTxShares(hash, namespace)
		if err == nil {
			return pointer, uint8(i), nil
		}
	}
	return nil, 0, fmt.Errorf("tx not found in any bundle")
}

func BundlesToShares(bundles []*Bundle, namespace string) []shares.Share {
	ss := []shares.Share{}
	for _, bundle := range bundles {
		s, _ := bundle.Shares(namespace)
		ss = append(ss, s...)
	}
	return ss
}

func GetSharesRoot(bundles []*Bundle, namespace string) []byte {
	ss := BundlesToShares(bundles, namespace)
	return merkle.HashFromByteSlices(shares.ToBytes(ss))
}

func GetSharesProofs(sp SharePointer, bundles []*Bundle, bundleNum int, ns string) []*merkle.Proof {
	offset := 0
	ss := []shares.Share{}

	for i := 0; i < len(bundles); i++ {
		s, _ := bundles[i].Shares(ns)
		ss = append(ss, s...)

		if bundleNum < i {
			offset += len(s)
		}
	}
	_, proofs := merkle.ProofsFromByteSlices(shares.ToBytes(ss))

	// adjust the index of the proof
	startProof := offset + sp.StartShare
	endProof := offset + sp.EndShare() + 1
	return proofs[startProof:endProof]
}
