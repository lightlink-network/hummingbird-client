package utils

import (
	"context"
	"math/big"

	"hummingbird/node/lightlink/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

func HashHeaderWithoutExtraData(header *ethtypes.Header) common.Hash {
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

func ToL2HeaderJson(header *ethtypes.Header) *L2HeaderJson {

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

type TxJson struct {
	ChainID  *big.Int       `json:"chainId,omitempty"`
	Nonce    uint64         `json:"nonce"`
	GasPrice *big.Int       `json:"gasPrice"`
	Gas      uint64         `json:"gas"`
	To       common.Address `json:"to"`
	Value    *big.Int       `json:"value"`
	Data     []byte         `json:"data"`
	V        uint8          `json:"v"`
	R        *big.Int       `json:"r"`
	S        *big.Int       `json:"s"`
}

func ToTxJson(tx *types.Transaction) *TxJson {
	to := common.Address{}
	if tx.To() != nil {
		to = *tx.To()
	}

	var chainID *big.Int
	if tx.ChainId().Uint64() != 0 {
		chainID = tx.ChainId()
	}

	v, r, s := tx.RawSignatureValues()
	return &TxJson{
		ChainID:  chainID,
		Nonce:    tx.Nonce(),
		GasPrice: tx.GasPrice(),
		Gas:      tx.Gas(),
		To:       to,
		Value:    tx.Value(),
		Data:     tx.Data(),
		V:        uint8(v.Uint64()),
		R:        r,
		S:        s,
	}
}
