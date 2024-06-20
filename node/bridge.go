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

	lastBlock, err := n.LightLink.GetBlock(rollupHeader.L2Height - 1)
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

func (n *Node) GetWithdrawalProof(withdrawalHash common.Hash) ([][]byte, error) {
	// TODO
	panic("not implemented")
}
