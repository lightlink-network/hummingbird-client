package node

import (
	"bytes"
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func randAddr() common.Address {
	buf := make([]byte, 20)
	crand.Read(buf)
	return common.BytesToAddress(buf)
}

func newRandomBundle(size int, withTxns bool) *Bundle {
	b := &Bundle{
		Blocks: make([]*types.Block, size),
	}

	for i := 0; i < size; i++ {
		prevHash := common.Hash{}
		if i > 0 {
			prevHash = b.Blocks[i-1].Hash()
		}

		b.Blocks[i] = types.NewBlockWithHeader(&types.Header{
			ParentHash: prevHash,
			Time:       uint64(i),
			Coinbase:   randAddr(),
			Number:     big.NewInt(int64(i)),
		})

		if withTxns {
			randN := 1 + rand.Intn(10)
			txns := make(types.Transactions, randN)
			for j := 0; j < randN; j++ {
				txns[j] = types.NewTransaction(uint64(j), b.Blocks[i].Coinbase(), big.NewInt(0), 100000, big.NewInt(0), nil)
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
