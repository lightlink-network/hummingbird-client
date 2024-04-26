package utils

import (
	"context"
	"math/big"

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

func HashHeaderWithoutExtraData(header *types.Header) common.Hash {
	header.Extra = common.Hex2Bytes("0x")
	return header.Hash()
}

type L2HeaderJson struct {
	ParentHash       common.Hash    `json:"parentHash"`
	UncleHash        common.Hash    `json:"uncleHash"`
	Beneficiary      common.Address `json:"beneficiary"`
	StateRoot        common.Hash    `json:"stateRoot"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	Difficulty       *big.Int       `json:"difficulty"`
	LogsBloom        []byte         `json:"logsBloom"`
	Number           *big.Int       `json:"number"`
	GasLimit         *big.Int       `json:"gasLimit"`
	GasUsed          *big.Int       `json:"gasUsed"`
	Timestamp        *big.Int       `json:"timestamp"`
	ExtraData        []byte         `json:"extraData"`
	MixHash          common.Hash    `json:"mixHash"`
	Nonce            []byte         `json:"nonce"`
}

func ToL2HeaderJson(header *types.Header) *L2HeaderJson {

	return &L2HeaderJson{
		ParentHash:       header.ParentHash,
		UncleHash:        header.UncleHash,
		Beneficiary:      header.Coinbase,
		StateRoot:        header.Root,
		TransactionsRoot: header.TxHash,
		ReceiptsRoot:     header.ReceiptHash,
		Difficulty:       header.Difficulty,
		LogsBloom:        header.Bloom[:],
		Number:           header.Number,
		GasLimit:         big.NewInt(int64(header.GasLimit)),
		GasUsed:          big.NewInt(int64(header.GasUsed)),
		Timestamp:        big.NewInt(int64(header.Time)),
		ExtraData:        common.Hex2Bytes("0x"),
		MixHash:          header.MixDigest,
		Nonce:            header.Nonce[:],
	}
}
