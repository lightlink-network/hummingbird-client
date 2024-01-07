package cmd

import (
	"hummingbird/node"
	"hummingbird/rollup"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var RollupCmd = &cobra.Command{
	Use:   "rollup",
	Short: "rollup client will create and submit new rollup blocks to layer 1",

	Run: func(cmd *cobra.Command, args []string) {
		cfg := LoadConfig()
		logger := DefaultLogger()
		ethKey := getEthKey()

		eth, err := node.NewEthereumRPC(node.EthereumRPCOpts{
			Endpoint:                   cfg.Ethereum.Endpoint,
			CanonicalStateChainAddress: common.HexToAddress(cfg.Ethereum.CanonicalStateChain),
			Signer:                     ethKey,
		})
		must(err)
		logger.Debug("Ethereum client successfully initialized")

		cel, err := node.NewCelestiaAPI(node.CelestiaAPIOpts{
			Endpoint:      cfg.Celestia.Endpoint,
			Token:         cfg.Celestia.Token,
			GRPC:          cfg.Celestia.GRPC,
			TendermintRPC: cfg.Celestia.TendermintRPC,
			Namespace:     cfg.Celestia.Namespace,
		})
		must(err)
		logger.Debug("Celestia client successfully initialized")

		ll, err := node.NewLightLinkClient(&node.LightLinkClientOpts{
			Endpoint: cfg.LightLink.Endpoint,
			Delay:    time.Duration(cfg.LightLink.Delay) * time.Millisecond,
		})
		must(err)
		logger.Debug("LightLink client successfully initialized")

		h, err := ll.GetHeight()
		must(err)
		logger.Info("LightLink height", "height", h)

		b, err := ll.GetBlock(h - 1)
		must(err)
		logger.Info("LightLink block", "txs", len(b.Transactions()))

		n := node.Node{
			Ethereum:  eth,
			Celestia:  cel,
			LightLink: ll,
		}

		_ = rollup.NewRollup(&n, &rollup.Opts{
			PollDelay:  time.Duration(cfg.Rollup.PollDelay) * time.Millisecond,
			BundleSize: cfg.Rollup.BundleSize,
			Logger:     logger.With("ctx", "Rollup"),
		})

		logger.Info("Rollup client successfully initialized")

	},
}
