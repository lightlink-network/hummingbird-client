package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func IsContract(client *ethclient.Client, address common.Address) (bool, error) {
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return false, err
	}

	return len(code) > 0, nil
}

func HashWithoutExtraData(block *types.Block) common.Hash {
	header := block.Header()
	header.Extra = common.Hex2Bytes("0x")
	return header.Hash()
}
