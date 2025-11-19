package utils

import (
	chainoracle "hummingbird/node/contracts/ChainOracle.sol"
	challenge "hummingbird/node/contracts/Challenge.sol"
	"math/big"

	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func CalculateMerkleRoot(leafs ...common.Hash) common.Hash {
	if len(leafs) == 0 {
		return common.Hash{}
	}

	var branches []common.Hash
	for i := 0; i < len(leafs); i += 2 {
		if i+1 > len(leafs)-1 {
			branch := crypto.Keccak256(leafs[i].Bytes(), leafs[i].Bytes())
			branches = append(branches, common.BytesToHash(branch))
		} else {
			branch := crypto.Keccak256(leafs[i].Bytes(), leafs[i+1].Bytes())
			branches = append(branches, common.BytesToHash(branch))
		}
	}

	if len(branches) != 1 {
		return CalculateMerkleRoot(branches...)
	}

	return branches[0]
}

func ToBinaryMerkleProof(proofs []*merkle.Proof) []chainoracle.BinaryMerkleProof {
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

func ToChallengeBinaryMerkleProof(proofs []*merkle.Proof) []challenge.BinaryMerkleProof {
	rowProofs := make([]challenge.BinaryMerkleProof, len(proofs))
	for i, proof := range proofs {
		sideNodes := make([][32]byte, len(proof.Aunts))
		for j, sideNode := range proof.Aunts {
			var bzSideNode [32]byte
			for k, b := range sideNode {
				bzSideNode[k] = b
			}
			sideNodes[j] = bzSideNode
		}
		rowProofs[i] = challenge.BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       big.NewInt(proof.Index),
			NumLeaves: big.NewInt(proof.Total),
		}
	}
	return rowProofs
}
