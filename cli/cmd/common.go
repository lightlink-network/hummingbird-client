package cmd

import (
	"crypto/ecdsa"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/tint"
)

func ConsoleLogger() *slog.Logger {
	w := os.Stderr
	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	}))

	return logger
}

func JSONLogger(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

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
