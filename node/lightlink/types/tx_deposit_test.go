package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// NewDepositTx creates an unsigned deposit transaction.
func TestNewDepositTx(t *testing.T) {
	tx := NewTx(&DepositTx{
		Nonce:    1,
		GasPrice: big.NewInt(1),
		Gas:      21000,
		To:       &common.Address{},
		Value:    big.NewInt(1),
		Data:     []byte{1},
	})

	// assert values
	assert.Equal(t, tx.Nonce(), uint64(1))
	assert.Equal(t, tx.GasPrice().Uint64(), uint64(1))
	assert.Equal(t, tx.Gas(), uint64(21000))
	assert.Equal(t, tx.To(), &common.Address{})
	assert.Equal(t, tx.Value().Uint64(), uint64(1))
	assert.Equal(t, tx.Data(), []byte{1})
}

// Test that the transaction type is set correctly.
func TestDepositTxType(t *testing.T) {
	tx := NewTx(&DepositTx{})
	assert.Equal(t, tx.Type(), uint8(DepositTxType))
}

// Test sign and recover.
func TestDepositTxSignAndRecover(t *testing.T) {
	// create test signer
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// create a signed transaction
	tx, err := SignNewTx(privateKey, NewLightLinkSigner(big.NewInt(88)), &DepositTx{
		Nonce:    1,
		GasPrice: big.NewInt(1),
		Gas:      21000,
		To:       &common.Address{},
		Value:    big.NewInt(1),
		Data:     []byte{1},
	})
	assert.NoError(t, err)

	// recover the sender
	sender, err := Sender(NewLightLinkSigner(big.NewInt(88)), tx)
	assert.NoError(t, err)

	// assert values
	assert.Equal(t, sender, crypto.PubkeyToAddress(privateKey.PublicKey))
}

// Test MarshalBinary and UnmarshalBinary.
func TestDepositTxMarshalBinary(t *testing.T) {
	// create a signed transaction
	addr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	tx := NewTx(&DepositTx{
		ChainID:  big.NewInt(88),
		Nonce:    8,
		GasPrice: big.NewInt(1000000000),
		Gas:      21000,
		To:       &addr,
		Value:    big.NewInt(0),
		Data:     []byte{0x4},
	})

	// net private key from hex
	prv, _ := crypto.HexToECDSA("8a0e7cb61f25b74a0b86cdc39e4e4a4d05b322ec4807cbaee50522a25aee6cd6")

	txSigned, _ := SignTx(tx, NewLightLinkSigner(big.NewInt(88)), prv)

	// marshal the transaction
	data, err := txSigned.MarshalBinary()
	assert.NoError(t, err)

	fmt.Println("0x" + common.Bytes2Hex(data))

	// unmarshal the transaction
	tx2 := &Transaction{}
	err = tx2.UnmarshalBinary(data)
	assert.NoError(t, err)

	// assert values
	assert.Equal(t, tx2.Nonce(), uint64(8))
	assert.Equal(t, tx2.GasPrice().Uint64(), uint64(1000000000))
	assert.Equal(t, tx2.Gas(), uint64(21000))
	assert.Equal(t, tx2.To(), &common.Address{})
	assert.Equal(t, tx2.Value().Uint64(), uint64(0))
	assert.Equal(t, tx2.Data(), []byte{4})
	assert.Equal(t, tx2.ChainId(), big.NewInt(88))
}
