package main

import (
	"hummingbird/cli/cmd"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/utils"
	"log/slog"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/spf13/viper"
)

func main() {
	viper.Set("config-path", "./ganache")
	viper.Set("log-level", "debug")
	cfg := config.Load()
	log := cmd.ConsoleLogger()
	ethKey := getEthKey()

	n, err := node.NewFromConfig(cfg, log, ethKey)
	panicIfError(err, "failed to create node from config")

	l2Bundle, err := fetchBundleFromL2(n)
	panicIfError(err, "failed to fetch bundle from L2")

	l2Blob, err := l2Bundle.Blob(n.Namespace())
	panicIfError(err, "failed to get blob from bundle")

	l2Shares, err := utils.BlobToShares(l2Blob)
	panicIfError(err, "failed to get shares from blob")

	log.Info("L2 Shares", "count", len(l2Shares))
	encoded := []byte{}
	for i, share := range l2Shares {
		d, _ := share.RawData()
		encoded = append(encoded, d...)
		log.Info("L2 Share", "i", i, "share", crypto.Keccak256Hash(d))
	}

	celestiaBundle, err := fetchBundleFromCelestia(n, log)
	panicIfError(err, "failed to fetch bundle from celestia")

	log.Info("Bundle", "blocks", len(celestiaBundle.Blocks))
}

func fetchBundleFromL2(n *node.Node) (*node.Bundle, error) {
	blocks, err := n.LightLink.GetBlocks(62207259, 62207259+10)
	if err != nil {
		return nil, err
	}

	return &node.Bundle{
		Blocks: blocks,
	}, nil
}

func fetchBundleFromCelestia(n *node.Node, log *slog.Logger) (*node.Bundle, error) {
	shares, err := n.Celestia.GetShares(&node.CelestiaPointer{
		Height:     1018885,
		ShareStart: 1,
		ShareLen:   13,
	})
	if err != nil {
		return nil, err
	}
	log.Info("DA Shares", "count", len(shares))

	data := make([]byte, 0)
	for i, share := range shares {
		d, _ := share.RawData()
		data = append(data, d...)
		log.Info("DA Share", "i", i, "share", crypto.Keccak256Hash(d))
	}

	b := &node.Bundle{}
	if err := b.DecodeRLP(data); err != nil {
		return nil, err
	}

	return b, nil
}
