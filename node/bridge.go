package node

import (
	lightlinkportal "hummingbird/node/contracts/LightLinkPortal.sol"

	"github.com/ethereum/go-ethereum/common"
)

func (n *Node) GenOutputProofV0(rblockHash common.Hash) (*lightlinkportal.TypesOutputRootProof, error) {
	rollupHeader, err := n.Ethereum.GetRollupHeaderByHash(rblockHash)
	if err != nil {
		return nil, err
	}

	lastBlock, err := n.LightLink.GetBlock(rollupHeader.L2Height)
	if err != nil {
		return nil, err
	}

	output, err := n.LightLink.GetOutputV0(lastBlock.Header())
	if err != nil {
		return nil, err
	}

	return &lightlinkportal.TypesOutputRootProof{
		Version:                  output.Version(),
		StateRoot:                output.StateRoot,
		MessagePasserStorageRoot: output.MessagePasserStorageRoot,
		LatestBlockhash:          output.BlockHash,
	}, nil
}

func (n *Node) GetWithdrawalProof(rblockHash, withdrawalRoot, withdrawalHash common.Hash) ([][]byte, error) {
	rollupHeader, err := n.Ethereum.GetRollupHeaderByHash(rblockHash)
	if err != nil {
		return nil, err
	}

	l2Height := rollupHeader.L2Height
	rawProof, err := n.LightLink.GetProof(n.LightLink.WithdrawalAddress(l2Height), []string{withdrawalHash.Hex()}, l2Height)
	if err != nil {
		return nil, err
	}

	proof := [][]byte{}
	for _, p := range rawProof.StorageProof[0].Proof {
		proof = append(proof, common.FromHex(p))
	}

	return proof, nil
}
