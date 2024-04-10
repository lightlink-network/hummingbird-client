package cmd

import (
	"fmt"
	"hummingbird/challenger"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/utils"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ChallengerChallengedaCmd = &cobra.Command{
	Use:        "challenge-da",
	Short:      "challengeda will create a challenge to a blocks dataroot inclusion on celestia",
	ArgAliases: []string{"block", "pointerIndex"},
	Args:       cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()

		// is dry run enabled?
		dryRun, _ := cmd.Flags().GetBool("dry")
		cfg.DryRun = dryRun

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		c := challenger.NewChallenger(n, &challenger.Opts{
			Logger: logger.With("ctx", "Challenger"),
			DryRun: dryRun,
		})

		// get block index from flags
		blockIndexRaw := args[0]
		blockIndex, err := strconv.ParseUint(blockIndexRaw, 10, 64)
		if err != nil {
			logger.Error("Failed to parse block index", "err", err)
			panic(err)
		}

		pointerIndex, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Error("Failed to parse pointer index", "err", err)
			panic(err)
		}

		tx, blockHash, err := c.ChallengeDA(blockIndex, uint8(pointerIndex))
		if err != nil {
			logger.Error("Failed to challenge data availability", "err", err)
			panic(err)
		}

		fmt.Println("Challenged data availability with tx:", tx.Hash().Hex(), "gas used:", tx.Gas(), "gas price:", tx.GasPrice().Uint64(), "block hash", blockHash.Hex())
	},
}
