package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	uint256Type, _ = abi.NewType("uint256", "uint256", nil)
	uint64Type, _  = abi.NewType("uint64", "uint64", nil)
	bytes32Type, _ = abi.NewType("bytes32", "bytes32", nil)
	bytesType, _   = abi.NewType("bytes", "bytes", nil)
	addressType, _ = abi.NewType("address", "address", nil)
)

func HashCanonicalStateChainHeader(header *CanonicalStateChainHeader) (common.Hash, error) {
	args := abi.Arguments{
		{Type: uint64Type},  // epoch
		{Type: uint64Type},  // l2Height
		{Type: bytes32Type}, // prevHash
		{Type: bytes32Type}, // txRoot
		{Type: bytes32Type}, // blockRoot
		{Type: bytes32Type}, // stateRoot
		{Type: uint64Type},  // celestiaHeight
		{Type: bytes32Type}, // celestiaDataRoot
	}

	enc, err := args.Pack(header.Epoch, header.L2Height, header.PrevHash, header.TxRoot, header.BlockRoot, header.StateRoot, header.CelestiaHeight, header.CelestiaDataRoot)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(crypto.Keccak256(enc)), nil
}
