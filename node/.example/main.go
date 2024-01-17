package main

import (
	"bytes"
	"fmt"
	"hummingbird/node"
	"hummingbird/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	bundle := generateBundle(10)
	buf, _ := bundle.EncodeRLP()

	blob, _ := utils.BytesToBlob("test", buf)
	s, _ := utils.BlobToShares(blob)

	h6, _ := rlp.EncodeToBytes(bundle.Blocks[6].Header())

	// get raw index of the header
	rawStart := bytes.Index(buf, h6)
	rawEnd := rawStart + len(h6)

	// get the share pointers
	startShare, startIndex := utils.RawIndexToSharesIndex(rawStart, s)
	endShare, endIndex := utils.RawIndexToSharesIndex(rawEnd, s)

	fmt.Printf("raw pointer   → %d %d\n", rawStart, rawEnd)
	fmt.Printf("start pointer → %d:%d\n", startShare, startIndex)
	fmt.Printf("end pointer   → %d:%d\n", endShare, endIndex)

	// reconstruct the block header
	// get the hdata from the shares
	hdata := []byte{}
	for i := startShare; i <= endShare; i++ {
		rawData, _ := s[i].RawData()

		if i == startShare {
			hdata = append(hdata, rawData[startIndex:]...)
		} else if i == endShare {
			hdata = append(hdata, rawData[:endIndex]...)
		} else {
			hdata = append(hdata, rawData...)
		}
	}

	fmt.Printf("hdata (%d) → %v\n\n", len(hdata), hdata)
	fmt.Printf("h6    (%d) → %v\n\n", len(h6), h6)

	// reconstruct the block header
	dech6 := &ethtypes.Header{}
	err := rlp.DecodeBytes(hdata, dech6)
	if err != nil {
		panic(err)
	}

	fmt.Printf("dech6 → %v\n\n", dech6.Root[:])

	// get the header from the selectedShares

	// rlpEnd := rlpStart + len(h6)
	// fmt.Println("rlp:start", rlpStart, rlpEnd)
	//
	// dech6 := &ethtypes.Header{}
	// err := rlp.DecodeBytes(buf[rlpStart:rlpEnd], dech6)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("rlp:dech6", dech6.Root[:])
	//
	// b, _ := utils.BytesToBlob("test", buf)
	// s, _ := utils.BlobToShares(b)
	//
	// fmt.Printf("shares → %d\n", len(s))
	//
	// for i := 0; i < len(s); i++ {
	// 	idx := utils.ShareDataStart(s[i])
	// 	size := len(s[i].ToBytes())
	//
	// 	fmt.Printf("shares(%d) → %d %d\n", i, size, size-idx)
	// }
	//
	// // step 1. get the list
	// _, list, _, _ := rlp.Split(buf)
	//
	// // step 2. read the elems
	// end := false
	// count := 0
	// for !end {
	// 	_, elem, rest, _ := rlp.Split(list)
	// 	list = rest
	// 	end = len(rest) < 1
	//
	// 	b := &ethtypes.Block{}
	// 	err := rlp.DecodeBytes(elem, b)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	fmt.Printf("elem %d → %v\n\n", count, b.Root())
	// 	count++
	// }

}

func generateBundle(size int) node.Bundle {
	blocks := []*ethtypes.Block{}
	prevHash := common.Hash{}

	for i := 0; i < size; i++ {
		tag := append([]byte{}, byte(i), byte(i), byte(i))

		h := &ethtypes.Header{
			ParentHash:  prevHash,
			UncleHash:   common.Hash{},
			Coinbase:    common.Address{},
			Root:        common.Hash(common.LeftPadBytes(tag, 32)),
			TxHash:      [32]byte{},
			ReceiptHash: [32]byte{},
			Bloom:       [256]byte{},
			Difficulty:  &big.Int{},
			Number:      &big.Int{},
			GasLimit:    uint64(i),
			GasUsed:     uint64(i),
			Time:        uint64(i),
			Extra:       []byte{},
			MixDigest:   [32]byte{},
			Nonce:       [8]byte{},
		}

		b := ethtypes.NewBlockWithHeader(h)
		blocks = append(blocks, b)
		prevHash = b.Hash()
	}

	return node.Bundle{blocks}
}

type SharePointer struct {
	StartShare int
	StartIndex int
	EndShare   int
	EndIndex   int
}
