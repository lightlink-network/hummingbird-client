package contracts

import (
	chainoracle "hummingbird/node/contracts/ChainOracle.sol"
	"math/big"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
)

func NewShareProof(proof *types.ShareProof, attestations chainoracle.AttestationProof) (*chainoracle.SharesProof, error) {
	return &chainoracle.SharesProof{
		Data:             proof.Data,
		ShareProofs:      toNamespaceMerkleMultiProofs(proof.ShareProofs),
		Namespace:        *namespace(proof.NamespaceID, uint8(proof.NamespaceVersion)),
		RowRoots:         toRowRoots(proof.RowProof.RowRoots),
		RowProofs:        toRowProofs(proof.RowProof.Proofs),
		AttestationProof: attestations,
	}, nil
}

// Methods for converting for use with DAVerifier Library
// See https://docs.celestia.org/developers/blobstream-proof-queries#converting-the-proofs-to-be-usable-in-the-daverifier-library

func toNamespaceMerkleMultiProofs(proofs []*tmproto.NMTProof) []chainoracle.NamespaceMerkleMultiproof {
	shareProofs := make([]chainoracle.NamespaceMerkleMultiproof, len(proofs))
	for i, proof := range proofs {
		sideNodes := make([]chainoracle.NamespaceNode, len(proof.Nodes))
		for j, node := range proof.Nodes {
			sideNodes[j] = *toNamespaceNode(node)
		}
		shareProofs[i] = chainoracle.NamespaceMerkleMultiproof{
			BeginKey:  big.NewInt(int64(proof.Start)),
			EndKey:    big.NewInt(int64(proof.End)),
			SideNodes: sideNodes,
		}
	}
	return shareProofs
}

func minNamespace(innerNode []byte) *chainoracle.Namespace {
	version := innerNode[0]
	var id [28]byte
	for i, b := range innerNode[1:29] {
		id[i] = b
	}
	return &chainoracle.Namespace{
		Version: [1]byte{version},
		Id:      id,
	}
}

func maxNamespace(innerNode []byte) *chainoracle.Namespace {
	version := innerNode[29]
	var id [28]byte
	for i, b := range innerNode[30:58] {
		id[i] = b
	}
	return &chainoracle.Namespace{
		Version: [1]byte{version},
		Id:      id,
	}
}

func toNamespaceNode(node []byte) *chainoracle.NamespaceNode {
	minNs := minNamespace(node)
	maxNs := maxNamespace(node)
	var digest [32]byte
	for i, b := range node[58:] {
		digest[i] = b
	}
	return &chainoracle.NamespaceNode{
		Min:    *minNs,
		Max:    *maxNs,
		Digest: digest,
	}
}

func namespace(namespaceID []byte, version uint8) *chainoracle.Namespace {
	var id [28]byte
	copy(id[:], namespaceID)
	return &chainoracle.Namespace{
		Version: [1]byte{version},
		Id:      id,
	}
}

func toRowRoots(roots []bytes.HexBytes) []chainoracle.NamespaceNode {
	rowRoots := make([]chainoracle.NamespaceNode, len(roots))
	for i, root := range roots {
		rowRoots[i] = *toNamespaceNode(root.Bytes())
	}
	return rowRoots
}

func toRowProofs(proofs []*merkle.Proof) []chainoracle.BinaryMerkleProof {
	rowProofs := make([]chainoracle.BinaryMerkleProof, len(proofs))
	for i, proof := range proofs {
		sideNodes := make([][32]byte, len(proof.Aunts))
		for j, sideNode := range proof.Aunts {
			var bzSideNode [32]byte
			for k, b := range sideNode {
				bzSideNode[k] = b
			}
			sideNodes[j] = bzSideNode
		}
		rowProofs[i] = chainoracle.BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       big.NewInt(proof.Index),
			NumLeaves: big.NewInt(proof.Total),
		}
	}
	return rowProofs
}

func toAttestationProof(
	nonce uint64,
	height uint64,
	blockDataRoot [32]byte,
	dataRootInclusionProof merkle.Proof,
) chainoracle.AttestationProof {
	sideNodes := make([][32]byte, len(dataRootInclusionProof.Aunts))
	for i, sideNode := range dataRootInclusionProof.Aunts {
		var bzSideNode [32]byte
		for k, b := range sideNode {
			bzSideNode[k] = b
		}
		sideNodes[i] = bzSideNode
	}

	return chainoracle.AttestationProof{
		TupleRootNonce: big.NewInt(int64(nonce)),
		Tuple: chainoracle.DataRootTuple{
			Height:   big.NewInt(int64(height)),
			DataRoot: blockDataRoot,
		},
		Proof: chainoracle.BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       big.NewInt(dataRootInclusionProof.Index),
			NumLeaves: big.NewInt(dataRootInclusionProof.Total),
		},
	}
}
