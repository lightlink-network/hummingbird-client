package contracts

import (
	chainoracle "hummingbird/node/contracts/ChainOracle.sol"
	"math/big"

	"github.com/tendermint/tendermint/types"
)

func NewShareProof(proof *types.ShareProof, attestations chainoracle.AttestationProof) (*chainoracle.SharesProof, error) {

	// TODO: is this correct????

	ns := chainoracle.Namespace{
		Version: [1]byte(proof.NamespaceID),
		Id:      [28]byte(proof.NamespaceID[:28]),
	}

	sp := []chainoracle.NamespaceMerkleMultiproof{}
	for _, s := range proof.ShareProofs {
		sideNodes := []chainoracle.NamespaceNode{}
		for _, sn := range s.Nodes {
			sideNodes = append(sideNodes, chainoracle.NamespaceNode{
				Min:    ns,           // Not sure if this is correct???
				Max:    ns,           // ?
				Digest: [32]byte(sn), // ?
			})
		}

		sp = append(sp, chainoracle.NamespaceMerkleMultiproof{
			BeginKey:  big.NewInt(int64(s.Start)),
			EndKey:    big.NewInt(int64(s.End)),
			SideNodes: sideNodes,
		})
	}

	rr := []chainoracle.NamespaceNode{}
	for _, r := range proof.RowProof.RowRoots {

		rr = append(rr, chainoracle.NamespaceNode{
			Min:    ns,                  // Not sure if this is correct???
			Max:    ns,                  // ?
			Digest: [32]byte(r.Bytes()), // ?
		})

	}

	rp := []chainoracle.BinaryMerkleProof{}
	for _, r := range proof.RowProof.Proofs {

		sideNodes := [][32]byte{}
		for _, sn := range r.Aunts {
			sideNodes = append(sideNodes, [32]byte(sn))
		}

		rp = append(rp, chainoracle.BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       big.NewInt(r.Index),
			NumLeaves: big.NewInt(r.Total),
		})
	}

	return &chainoracle.SharesProof{
		Data:             proof.Data,
		ShareProofs:      sp,
		Namespace:        ns,
		RowRoots:         rr,
		RowProofs:        rp,
		AttestationProof: attestations,
	}, nil
}
