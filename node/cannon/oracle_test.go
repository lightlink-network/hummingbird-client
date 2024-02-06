package cannon_test

import (
	"hummingbird/node/cannon"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/stretchr/testify/assert"
)

func preImageLog(t *testing.T) func(cannon.PreImages) {
	return func(p cannon.PreImages) {
		for k, v := range p {
			t.Logf("New Preimage: %s -> %x", k, v)
		}
	}
}

func setupOracle(t *testing.T) *cannon.Oracle {
	publicEthUrl := "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	eth, err := ethclient.Dial(publicEthUrl)
	assert.NoError(t, err)
	return cannon.NewOracle(gethclient.New(eth.Client()), eth, nil)
}

func TestOracle_FetchAccount(t *testing.T) {
	oracle := setupOracle(t)

	err := oracle.PreFetchAccount(big.NewInt(19150042), common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), preImageLog(t))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(oracle.PreImages()), 1)
}

func TestOracle_FetchStorage(t *testing.T) {
	oracle := setupOracle(t)

	err := oracle.PreFetchStorage(big.NewInt(19150042), common.HexToAddress("0x6f259637dcd74c767781e37bc6133cd6a68aa161"), common.HexToHash("0x0"), preImageLog(t))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(oracle.PreImages()), 1)
}

func TestOracle_PreFetchBlock(t *testing.T) {
	oracle := setupOracle(t)

	err := oracle.PreFetchBlock(big.NewInt(19150042), true)
	assert.NoError(t, err)
	// this is first block only first input
	assert.GreaterOrEqual(t, len(oracle.PreImages()), 1)
	assert.NotEmpty(t, oracle.Inputs()[0].Big().Bytes())
	assert.Empty(t, oracle.Inputs()[1].Big().Bytes())
	assert.Empty(t, oracle.Outputs()[0].Big().Bytes())
	// TODO check hash is input 0 – but we need to get the hash first

	err = oracle.PreFetchBlock(big.NewInt(19150043), false)
	assert.NoError(t, err)
	// this is second block
	assert.GreaterOrEqual(t, len(oracle.PreImages()), 1)
	assert.NotEmpty(t, oracle.Inputs()[1].Big().Bytes())
	assert.NotEmpty(t, oracle.Outputs()[0].Big().Bytes())
}
