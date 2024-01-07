package main

import (
	"hummingbird/cli/cmd"
	"hummingbird/node"
	"time"
)

func main() {
	cfg := cmd.LoadConfig()
	logger := cmd.DefaultLogger()

	ll, err := node.NewLightLinkClient(&node.LightLinkClientOpts{
		Endpoint: cfg.LightLink.Endpoint,
		Delay:    time.Duration(cfg.LightLink.Delay) * time.Millisecond,
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Failed to connect to LightLink", "error", err)
		panic(err)
	}

	h, err := ll.GetHeight()
	if err != nil {
		logger.Error("Failed to get LightLink height", "error", err)
		panic(err)
	}

	logger.Info("LightLink height", "height", h)

	b, err := ll.GetBlock(62050067)
	if err != nil {
		logger.Error("Failed to get LightLink block", "error", err)
		panic(err)
	}

	logger.Info("LightLink block", "num", b.Number(), "hash", b.Hash().Hex(), "txs", len(b.Transactions()))
}
