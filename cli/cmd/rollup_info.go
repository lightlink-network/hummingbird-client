package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/rollup"
	"hummingbird/utils"
	"time"

	"github.com/spf13/cobra"
)

var RollupInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "info will print information about the current rollup state",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := ConsoleLogger()
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		r := rollup.NewRollup(n, &rollup.Opts{
			PollDelay:             time.Duration(cfg.Rollup.PollDelay) * time.Millisecond,
			BundleSize:            cfg.Rollup.BundleSize,
			StoreCelestiaPointers: cfg.Rollup.StoreCelestiaPointers,
			StoreHeaders:          cfg.Rollup.StoreHeaders,
			Logger:                logger.With("ctx", "Rollup"),
		})

		info, err := r.GetInfo()
		utils.NoErr(err)

		// fmt.Println(" ")
		// fmt.Println("Rollup Height:", info.RollupHeight)
		// fmt.Println("L2 Blocks Rolled Up:", info.L2BlocksRolledUp)
		// fmt.Println("L2 Blocks Todo:", info.L2BlocksTodo)
		// fmt.Println(" ")
		// fmt.Println("Latest Rollup Block:")
		// fmt.Println("	Hash:", info.LatestRollup.Hash.Hex())
		// fmt.Println("	Bundle Size:", info.LatestRollup.BundleSize)
		// fmt.Println("	Epoch:", info.LatestRollup.Epoch)
		// fmt.Println("	Height:", info.LatestRollup.L2Height)
		// fmt.Println("	Prev Hash:", common.Hash(info.LatestRollup.PrevHash).Hex())
		// fmt.Println("	StateRoot:", common.Hash(info.LatestRollup.StateRoot).Hex())
		// fmt.Println("	BlockRoot:", common.Hash(info.LatestRollup.BlockRoot).Hex())
		// fmt.Println("	TxRoot:", common.Hash(info.LatestRollup.TxRoot).Hex())
		// fmt.Println("	Data Availability:")
		// fmt.Println("		Celestia Height:", info.LatestRollup.CelestiaHeight)
		// fmt.Println("		Celestia Data Root:", common.Hash(info.LatestRollup.CelestiaDataRoot).Hex())
		// fmt.Println(" ")

		fmt.Println(" ")
		fmt.Println(utils.MarshalText(info))
		fmt.Println(" ")
	},
}
