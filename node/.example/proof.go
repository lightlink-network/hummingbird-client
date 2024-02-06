package main

import (
	"hummingbird/node"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	n, log := setup()

	celPointer := &node.CelestiaPointer{
		Height:     1018885,
		ShareStart: 1,
		ShareLen:   13,
	}

	log.Info("Fetch Bundle from Celestia")
	s, err := n.Celestia.GetShares(celPointer)

	panicIfError(err, "failed to fetch shares from Celestia")
	bundle, err := node.NewBundleFromShares(s)
	panicIfError(err, "failed to create bundle from shares")
	log.Info("Downloaded Bundle", "blocks", len(bundle.Blocks), "shares", len(s))

	log.Info("Get SharePointer for the 3rd block in the bundle")
	sharePointer, err := bundle.FindHeaderShares(bundle.Blocks[2].Hash(), n.Namespace())
	panicIfError(err, "failed to get share pointer")

	log.Info("Check Shares contain the 3rd block in the bundle")
	header, err := sharesPointerToHeader(sharePointer, s[sharePointer.StartShare:])
	panicIfError(err, "failed to convert shares to header")

	if header.Hash() != bundle.Blocks[2].Hash() {
		panic("shares do not contain the 3rd block in the bundle")
	}

	log.Info("Get Share Proof for the 3rd block in the bundle")
	proof, err := n.Celestia.GetSharesProof(celPointer, sharePointer)
	panicIfError(err, "failed to get share proof")

	log.Info("Got proof", "proof", len(proof.Data))
	proofShares, err := bytesToShares(proof.Data)
	panicIfError(err, "failed to convert proof to shares")

	log.Info("Got proof shares", "shares", len(proofShares))

	log.Info("Get Header from proof shares")
	proofHeader, err := sharesPointerToHeader(sharePointer, proofShares)
	panicIfError(err, "failed to convert proof shares to header")

	if proofHeader.Hash() != bundle.Blocks[2].Hash() {
		panic("proof shares do not contain the 3rd block in the bundle")
	}

	// proofBundle, err := node.NewBundleFromShares(proofShares)
	// panicIfError(err, "failed to create bundle from proof shares")

	// log.Info("Got proof bundle", "blocks", len(proofBundle.Blocks), "shares", len(proofShares))
}

func bytesToShares(data [][]byte) ([]shares.Share, error) {
	s := []shares.Share{}

	for _, d := range data {
		x, err := shares.NewShare(d)
		if err != nil {
			return nil, err
		}

		s = append(s, *x)
	}

	return s, nil
}

func sharesPointerToHeader(pointer *node.SharePointer, s []shares.Share) (*ethtypes.Header, error) {
	data := []byte{}
	for i, r := range pointer.Ranges {
		data = append(data, s[i].ToBytes()[r.Start:r.End]...)
	}

	header := &ethtypes.Header{}
	return header, rlp.DecodeBytes(data, &header)
}
