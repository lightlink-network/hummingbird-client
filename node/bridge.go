package node

import (
	"fmt"
	lightlinkportal "hummingbird/node/contracts/LightLinkPortal.sol"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
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
	slot, err := getSlot(withdrawalHash)
	if err != nil {
		return nil, err
	}

	rollupHeader, err := n.Ethereum.GetRollupHeaderByHash(rblockHash)
	if err != nil {
		return nil, err
	}

	l2Height := rollupHeader.L2Height
	rawProof, err := n.LightLink.GetProof(n.LightLink.WithdrawalAddress(l2Height), []string{slot.Hex()}, l2Height)
	if err != nil {
		return nil, err
	}

	proof := [][]byte{}
	for _, p := range rawProof.StorageProof[0].Proof {
		proof = append(proof, common.FromHex(p))
	}

	proof = fixProof(crypto.Keccak256Hash(slot.Bytes()).Hex(), proof)
	return proof, nil
}

// see https://github.com/ethereum-optimism/optimism/blob/f8143c8cbc4cc0c83922c53f17a1e47280673485/packages/sdk/src/utils/message-utils.ts#L42
func getSlot(messageHash common.Hash) (common.Hash, error) {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	arguments := abi.Arguments{{Type: bytes32Type}, {Type: uint256Type}}

	// Convert messageHash to bytes32
	msgHashBytes := messageHash.Bytes()    // key
	zeroHashBytes := common.Hash{}.Bytes() // 0 slot

	// Encode the arguments
	encodedData, err := arguments.Pack(msgHashBytes, zeroHashBytes)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to encode data: %w", err)
	}

	return crypto.Keccak256Hash(encodedData), nil
}

// Fix for the case where the final proof element is less than 32 bytes and the element exists inside of a branch node
// see: https://github.com/ethereum-optimism/optimism/blob/f8143c8cbc4cc0c83922c53f17a1e47280673485/packages/sdk/src/utils/merkle-utils.ts#L57
func fixProof(key string, proof [][]byte) [][]byte {

	// get the last element of the proof
	last := proof[len(proof)-1]
	var lastDecoded []any
	rlp.DecodeBytes(last, &lastDecoded)

	if len(lastDecoded) != 17 {
		return proof
	}

	for _, item := range lastDecoded {
		if node, ok := item.([]any); ok {

			// ???
			// const suffix = toHexString(item[0]).slice(3)
			// if (key.endsWith(suffix)) {
			// 	modifiedProof.push(toHexString(rlp.encode(item)))
			//}

			suffix := hexutil.Encode(node[0].([]byte))[2:]
			if strings.HasSuffix(key, suffix) {
				buf, _ := rlp.EncodeToBytes(item)
				proof = append(proof, buf)
			}
		}
	}

	return proof
}
