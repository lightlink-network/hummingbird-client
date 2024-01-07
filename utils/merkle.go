package utils

import (
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
