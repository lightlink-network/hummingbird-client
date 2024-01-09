package cmd

import (
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/rollup"
	"hummingbird/utils"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RollupStartCmd.Flags().Bool("dry", false, "dry run will not submit the rollup block to the L1 rollup contract, and will not upload real data to celestia")
}

var RollupStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start will start the rollup node",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := ConsoleLogger()
		ethKey := getEthKey()

		// is dry run enabled?
		dryRun, _ := cmd.Flags().GetBool("dry")
		cfg.DryRun = dryRun

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		r := rollup.NewRollup(n, &rollup.Opts{
			L1PollDelay:           time.Duration(cfg.Rollup.L1PollDelay) * time.Millisecond,
			L2PollDelay:           time.Duration(cfg.Rollup.L2PollDelay) * time.Millisecond,
			BundleSize:            cfg.Rollup.BundleSize,
			StoreCelestiaPointers: cfg.Rollup.StoreCelestiaPointers,
			StoreHeaders:          cfg.Rollup.StoreHeaders,
			Logger:                logger.With("ctx", "Rollup"),
			DryRun:                dryRun,
		})

		// If dry run is enabled, swap out celestia with a mock celestia client.
		if dryRun {
			logger.Warn("DryRun is enabled, using mock celestia client")
			r.Celestia = node.NewCelestiaMock(cfg.Celestia.Namespace)
		}

		for {
			err = r.Run()
			if err != nil {
				logger.Error("Rollup.Run failed", "err", err, "retry_in", "5s")
			}
			time.Sleep(5 * time.Second)
		}

	},
}
