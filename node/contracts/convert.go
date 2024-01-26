package contracts

import (
	"fmt"
	chainloader "hummingbird/node/contracts/ChainLoader.sol"

	"github.com/tendermint/tendermint/types"
)

func NewShareProof(proof *types.ShareProof) (*chainloader.SharesProof, error) {
	return nil, fmt.Errorf("not implemented")

	// sp := []chainloader.NamespaceMerkleMultiproof{}
	// for _, s := range proof.ShareProofs {

	// 	sp = append(sp, chainloader.NamespaceMerkleMultiproof{
	// 		BeginKey:  big.NewInt(int64(s.Start)),
	// 		EndKey:    big.NewInt(int64(s.End)),
	// 		SideNodes: []chainloader.NamespaceNode{},
	// 	})
	// }

	// ns := chainloader.Namespace{
	// 	Version: [1]byte(proof.NamespaceID),
	// 	Id:      [28]byte(proof.NamespaceID[:28]),
	// }

	// rr := []chainloader.NamespaceNode{}
	// for _, r := range proof.RowProof.RowRoots {

	// 	rr = append(rr, chainloader.NamespaceNode{
	// 		Min:    ns,
	// 		Max:    ns,
	// 		Digest: [32]byte(r.Bytes()),
	// 	})
	// }

	// rp := []chainloader.BinaryMerkleProof{}
	// for _, r := range proof.RowProof.Proofs {

	// 	rp = append(rp, chainloader.BinaryMerkleProof{
	// 		SideNodes: [][32]byte{},
	// 		Key:       big.NewInt(r.Index),
	// 		NumLeaves: big.NewInt(r.Total),
	// 	})
	// }

	// return &chainloader.SharesProof{
	// 	Data:             proof.Data,
	// 	ShareProofs:      sp,
	// 	Namespace:        ns,
	// 	RowRoots:         rr,
	// 	RowProofs:        []chainloader.BinaryMerkleProof{},
	// 	AttestationProof: chainloader.AttestationProof{

	// 	},
	// }, nil
}
