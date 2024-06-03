package cmd

import (
	"fmt"
	hbcmd "hummingbird/cli/hb/cmd"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/rollup"
	"hummingbird/utils"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	PushBadCmd = &cobra.Command{
		Use:   "push-bad [reasons...]",
		Short: "push-bad will push a bad block to Layer 1",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("push-bad called")
			cfg := config.Load()
			logger := hbcmd.GetLogger(viper.GetString("log-type"))
			ethKey := getEthKey()

			logger.Info("Pushing bad block", "reasons", args)

			n, err := node.NewFromConfig(cfg, logger, ethKey)
			utils.NoErr(err)

			r := rollup.NewRollup(n, &rollup.Opts{
				L1PollDelay: time.Duration(cfg.Rollup.L1PollDelay) * time.Millisecond,
				L2PollDelay: time.Duration(cfg.Rollup.L2PollDelay) * time.Millisecond,
				BundleSize:  100,
				BundleCount: 1,
				Store:       cfg.Rollup.Store,
				Logger:      logger.With("ctx", "Rollup"),
				DryRun:      false,
			})

			// get reasons
			reasons := strings.Join(args, " ")
			if !strings.Contains(reasons, "epoch") && !strings.Contains(reasons, "l2height") {
				logger.Error("Invalid reasons", "reasons", reasons)
				logger.Info("Tip: use 'epoch' or 'l2height' as reasons")
				return
			}

			// push bad block
			logger.Info("Creating bad block...")
			b, err := r.CreateNextBlock()
			if err != nil {
				logger.Error("Failed to create bad block", "err", err)
				return
			}
			logger.Info("Created bad block!", "bundles", len(b.Bundles))

			logger.Info("Getting rollup head...")
			head, err := r.GetRollupHead()
			if err != nil {
				logger.Error("Failed to get rollup head", "err", err)
				return
			}

			// distort the block
			if strings.Contains(reasons, "epoch") {
				b.Epoch = head.Epoch - 1
			}
			if strings.Contains(reasons, "l2height") {
				b.L2Height = head.L2Height - 1
			}

			logger.Info("Submitting bad block...", "epoch", b.Epoch, "l2height", b.L2Height)
			tx, err := r.SubmitBlock(b)
			if err != nil {
				logger.Error("Failed to submit bad block", "err", err)
				return
			}

			logger.Info("Submitted bad block", "tx", tx.Hash().String())
		},
	}
)

func init() {
	RootCmd.AddCommand(PushBadCmd)
}
