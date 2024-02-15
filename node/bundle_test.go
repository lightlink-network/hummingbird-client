package node

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func randAddr() common.Address {
	buf := make([]byte, 20)
	crand.Read(buf)
	return common.BytesToAddress(buf)
}

func randPrivKey() *ecdsa.PrivateKey {
	k, _ := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	return k
}

func newRandomBundle(size int, withTxns bool) *Bundle {
	b := &Bundle{
		Blocks: make([]*ethtypes.Block, size),
	}

	for i := 0; i < size; i++ {
		prevHash := common.Hash{}
		if i > 0 {
			prevHash = b.Blocks[i-1].Hash()
		}

		b.Blocks[i] = ethtypes.NewBlockWithHeader(&ethtypes.Header{
			ParentHash: prevHash,
			Time:       uint64(i),
			Coinbase:   randAddr(),
			Number:     big.NewInt(int64(i)),
		})

		if withTxns {
			signer := ethtypes.NewEIP155Signer(big.NewInt(1))
			randN := 1 + rand.Intn(10)
			txns := make(ethtypes.Transactions, randN)
			for j := 0; j < randN; j++ {
				to := randAddr()
				rtx := ethtypes.NewTx(&ethtypes.LegacyTx{
					Nonce:    uint64(j),
					To:       &to,
					Value:    big.NewInt(10000),
					Gas:      21000,
					GasPrice: big.NewInt(10000),
					Data:     []byte{},
				})

				stx, _ := ethtypes.SignTx(rtx, signer, randPrivKey())
				txns[j] = stx
			}

			b.Blocks[i] = b.Blocks[i].WithBody(txns, nil)
		}
	}

	return b
}

func TestBundle_EncodeRLP(t *testing.T) {
	b := newRandomBundle(5, true)

	buf, err := b.EncodeRLP()
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}

func TestBundle_DecodeRLP(t *testing.T) {
	b := newRandomBundle(5, true)

	encoded, err := b.EncodeRLP()
	assert.NoError(t, err)

	decoded := &Bundle{}
	err = decoded.DecodeRLP(encoded)
	assert.NoError(t, err)

	// check that blocks are distinct
	assert.NotEqual(t, b.Blocks[0].Hash(), decoded.Blocks[4].Hash())

	// now check that decoded blocks match original blocks
	assert.Equal(t, b.Blocks[0].Hash(), decoded.Blocks[0].Hash())
	assert.Equal(t, b.Blocks[1].Hash(), decoded.Blocks[1].Hash())
	assert.Equal(t, b.Blocks[2].Hash(), decoded.Blocks[2].Hash())
	assert.Equal(t, b.Blocks[3].Hash(), decoded.Blocks[3].Hash())
	assert.Equal(t, b.Blocks[4].Hash(), decoded.Blocks[4].Hash())
}

func TestBundleTxInclusion(t *testing.T) {
	b := newRandomBundle(2, true)
	tx := b.Blocks[0].Transactions()[0]

	encTx, err := rlp.EncodeToBytes(tx)
	assert.NoError(t, err)

	encB, err := b.EncodeRLP()
	assert.NoError(t, err)
	// check if encTx bytes are included in encB bytes
	assert.True(t, bytes.Contains(encB, encTx))
}

func TestBundle_FindTxShares(t *testing.T) {
	b := newRandomBundle(5, true)

	// select a target tx
	tx := b.Blocks[2].Transactions()[0]
	shares, err := b.Shares("test")
	assert.NoError(t, err)

	// check tx can be encoded and decoded
	encTx, err := rlp.EncodeToBytes(tx)
	assert.NoError(t, err)

	decTx := &ethtypes.Transaction{}
	err = rlp.DecodeBytes(encTx, decTx)
	assert.NoError(t, err)
	assert.Equal(t, tx.Hash().Hex(), decTx.Hash().Hex())

	pointer, err := b.FindTxShares(tx.Hash(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, pointer)

	// check that the pointer is valid
	foundTx, err := sharesPointerToTx(pointer, shares)
	assert.NoError(t, err)
	assert.Equal(t, tx.Hash().Hex(), foundTx.Hash().Hex())
}

func sharesPointerToTx(pointer *SharePointer, s []shares.Share) (*ethtypes.Transaction, error) {
	data := []byte{}
	for i, r := range pointer.Ranges {
		data = append(data, s[i+pointer.StartShare].ToBytes()[r.Start:r.End]...)
	}

	tx := &ethtypes.Transaction{}
	err := rlp.DecodeBytes(data, &tx)
	return tx, err
}
