package main

import (
	"hummingbird/node"
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	n, log := setup()

	header, err := n.Ethereum.GetRollupHead()
	panicIfError(err, "failed to get rollup head")

	hash, err := contracts.HashCanonicalStateChainHeader(&header)
	panicIfError(err, "failed to hash header")

	pointer := &node.CelestiaPointer{
		Height:     header.CelestiaHeight,
		ShareStart: header.CelestiaShareStart,
		ShareLen:   header.CelestiaShareLen,
	}

	shares, err := n.Celestia.GetShares(pointer)
	panicIfError(err, "failed to get shares from Celestia")

	bundle, err := node.NewBundleFromShares(shares)
	panicIfError(err, "failed to create bundle from shares")

	log.Info("Downloaded Bundle", "blocks", len(bundle.Blocks), "shares", len(shares))

	sharePointer, err := bundle.FindHeaderShares(bundle.Blocks[2].Hash(), n.Namespace())
	panicIfError(err, "failed to get share pointer")

	shareProof, err := n.Celestia.GetSharesProof(pointer, sharePointer)
	panicIfError(err, "failed to get share proof")

	log.Info("Got share proof", "proof", len(shareProof.Data))

	// data := []byte{}
	// for i, r := range sharePointer.Ranges {
	// 	data = append(data, shareProof.Data[i][r.Start:r.End]...)
	// }
	// log.Info("Raw", "bytes", fmt.Sprintf("%x", data))

	// for i, r := range shareProof.Data {
	// 	log.Info("Proof", "i", i, "bytes", fmt.Sprintf("%x", r))
	// }

	// h := &ethtypes.Header{}
	// if err := rlp.DecodeBytes(data, h); err != nil {
	// 	panic(err)
	// }

	// 	log.Info("Decoded", "header", h)

	data := hexutil.MustDecode("0xf9021aa0ce095cb5cd4725f71278ce79cb4589e5a87147fcc148fdf587292a540ee15acca0000000000000000000000000000000000000000000000000000000000000000094dfad157b8d4e58c26bf9b947f8e75b5adbc7822ba03903de7f5290e9ef5974c2789c47778c69bff45299b10c2c2046774a6baec48fa00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008201f48403b5351d83e4e1c08084659bf868a00000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000088c9e41bfa0b90b3aa")
	h := &ethtypes.Header{}
	if err := rlp.DecodeBytes(data, h); err != nil {
		panic(err)
	}

	h.EncodeRLP()

	rlp.EncodeToBytes(h)

	log.Info("Decoded", "header", h)

	celProof, err := n.Celestia.GetProof(pointer)
	panicIfError(err, "failed to get proof")

	tx, err := n.Ethereum.ProvideShares(hash, shareProof, celProof)
	panicIfError(err, "failed to provide shares")
	n.Ethereum.Wait(tx.Hash())

	log.Info("Provided shares", "tx", tx.Hash().Hex())

	tx, err = n.Ethereum.ProvideHeader(hash, shareProof.Data, *sharePointer)
	panicIfError(err, "failed to provide header")

	log.Info("Provided header", "tx", tx.Hash().Hex())
}
