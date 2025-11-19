package main

import (
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	// Import the node package to access the types and functions
	"hummingbird/node"
	"hummingbird/node/lightlink/types"
)

func main() {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("Starting Celestia debug test")

	// Create a real Celestia client for debugging
	logger.Info("Creating Celestia client...")
	celestia, err := createCelestiaClient(logger)
	if err != nil {
		logger.Error("Failed to create Celestia client", "error", err)
		os.Exit(1)
	}
	logger.Info("Celestia client created successfully")

	// Create sample test blocks
	logger.Info("Creating test bundle...")
	bundle := createTestBundle()
	logger.Info("Test bundle created", "blocks", len(bundle.Blocks))

	// Test publishing the bundle
	logger.Info("Publishing test bundle to real Celestia network", "blocks", len(bundle.Blocks))

	pointer, gasPrice, err := celestia.PublishBundle(*bundle)
	if err != nil {
		logger.Error("Failed to publish bundle", "error", err)
		os.Exit(1)
	}

	logger.Info("Bundle published successfully to Celestia",
		"height", pointer.Height,
		"shareStart", pointer.ShareStart,
		"shareLen", pointer.ShareLen,
		"gasPrice", gasPrice,
		"txHash", pointer.TxHash.Hex())
}

// createCelestiaClient creates a real Celestia client with configuration
func createCelestiaClient(logger *slog.Logger) (node.Celestia, error) {
	// Configuration for Celestia connection using provided test details
	opts := node.CelestiaClientOpts{
		Endpoint:      "http://54.151.222.180:26658",                                                                                                                                                                                                                                         // Updated Celestia endpoint
		Token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdLCJOb25jZSI6IlYzczdtYVRXek4yU3U4Um5Ga0lSTzlBQk9QbWZMSForWDY2ekZRY3lIMEk9IiwiRXhwaXJlc0F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoifQ.yBLavfY7jwrSVZZlSGKZ-5kchMQvp9qpvzFeSoYQp0k", // Updated token
		TendermintRPC: "http://rpc-mocha.pops.one:26657",                                                                                                                                                                                                                                     // Updated Tendermint RPC endpoint
		Namespace:     "testdebug",                                                                                                                                                                                                                                                           // Updated namespace
		Logger:        logger,
		GasAPI:        "https://api-mocha-4.celenium.io/v1/gas/price", // Updated gas API
		GasPrice:      0.005,                                          // Same gas price
		Retries:       3,
		RetryDelay:    120 * time.Second, // Updated retry delay: 120000ms = 120s
	}
	return node.NewCelestiaClient(opts)
}

// createTestBundle creates a sample bundle with test blocks for publishing
func createTestBundle() *node.Bundle {
	// Create test header
	header := &ethtypes.Header{
		ParentHash:  common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Root:        common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef"),
		TxHash:      ethtypes.EmptyTxsHash,
		ReceiptHash: ethtypes.EmptyReceiptsHash,
		Bloom:       ethtypes.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(1),
		GasLimit:    21000,
		GasUsed:     21000,
		Time:        1234567890,
		Extra:       []byte("test"),
	}

	// Create test block
	block := types.NewBlockWithHeader(header)

	// Create bundle with the test block
	bundle := &node.Bundle{
		Blocks: []*types.Block{block},
	}

	return bundle
}
