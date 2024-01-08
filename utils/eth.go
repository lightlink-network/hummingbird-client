package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func IsContract(client *ethclient.Client, address common.Address) (bool, error) {
	code, err := client.CodeAt(nil, address, nil)
	if err != nil {
		return false, err
	}

	return len(code) > 0, nil
}
