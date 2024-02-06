package main

import (
	"crypto/ecdsa"
	"hummingbird/cli/cmd"
	"hummingbird/config"
	"hummingbird/node"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

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

func panicIfError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			panic(msg[0] + ": " + err.Error())
		}
		panic(err)
	}
}

func setup() (*node.Node, *slog.Logger) {
	viper.Set("config-path", "./ganache")
	viper.Set("log-level", "debug")
	cfg := config.Load()
	log := cmd.ConsoleLogger()
	ethKey := getEthKey()

	n, err := node.NewFromConfig(cfg, log, ethKey)
	panicIfError(err, "failed to create node from config")

	return n, log
}
