package cmd

import (
	"crypto/ecdsa"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func DefaultLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	})

	logger := slog.New(handler).With("app", "hummingbird")
	return logger
}

func getEthKey() *ecdsa.PrivateKey {
	key := os.Getenv("ETH_KEY")
	if key == "" {
		panic("env ETH_KEY not set")
	}

	ethKey, err := crypto.ToECDSA(hexutil.MustDecode(key))
	if err != nil {
		panic("Failed to decode ETH_KEY: " + err.Error())
	}

	return ethKey
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
