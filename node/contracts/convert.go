package contracts

import (
	chainloader "hummingbird/node/contracts/ChainLoader.sol"
	"math/big"

	"github.com/tendermint/tendermint/types"
)

func NewShareProof(proof *types.ShareProof, attestations chainloader.AttestationProof) (*chainloader.SharesProof, error) {

	// TODO: is this correct????

	ns := chainloader.Namespace{
		Version: [1]byte(proof.NamespaceID),
		Id:      [28]byte(proof.NamespaceID[:28]),
	}

	sp := []chainloader.NamespaceMerkleMultiproof{}
	for _, s := range proof.ShareProofs {
		sideNodes := []chainloader.NamespaceNode{}
		for _, sn := range s.Nodes {
			sideNodes = append(sideNodes, chainloader.NamespaceNode{
				Min:    ns,           // Not sure if this is correct???
				Max:    ns,           // ?
				Digest: [32]byte(sn), // ?
			})
		}

		sp = append(sp, chainloader.NamespaceMerkleMultiproof{
			BeginKey:  big.NewInt(int64(s.Start)),
			EndKey:    big.NewInt(int64(s.End)),
			SideNodes: sideNodes,
		})
	}

	rr := []chainloader.NamespaceNode{}
	for _, r := range proof.RowProof.RowRoots {

		rr = append(rr, chainloader.NamespaceNode{
			Min:    ns,                  // Not sure if this is correct???
			Max:    ns,                  // ?
			Digest: [32]byte(r.Bytes()), // ?
		})

	}

	rp := []chainloader.BinaryMerkleProof{}
	for _, r := range proof.RowProof.Proofs {

		sideNodes := [][32]byte{}
		for _, sn := range r.Aunts {
			sideNodes = append(sideNodes, [32]byte(sn))
		}

		rp = append(rp, chainloader.BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       big.NewInt(r.Index),
			NumLeaves: big.NewInt(r.Total),
		})
	}

	return &chainloader.SharesProof{
		Data:             proof.Data,
		ShareProofs:      sp,
		Namespace:        ns,
		RowRoots:         rr,
		RowProofs:        rp,
		AttestationProof: attestations,
	}, nil
}
