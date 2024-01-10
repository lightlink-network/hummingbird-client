package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	DefenderDefendDaCmd.Flags().Bool("dry", false, "dry run will not submit the rollup block to the L1 rollup contract, and will not upload real data to celestia")
	DefenderDefendDaCmd.Flags().String("tx", "", "celestia tx hash in which data was submitted")
	DefenderDefendDaCmd.MarkFlagRequired("tx")
}

var DefenderDefendDaCmd = &cobra.Command{
	Use:        "defend-da",
	Short:      "defend-da will defend against a data availability challenge",
	ArgAliases: []string{"block"},
	Args:       cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := ConsoleLogger()
		ethKey := getEthKey()

		// is dry run enabled?
		dryRun, _ := cmd.Flags().GetBool("dry")
		cfg.DryRun = dryRun

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
			DryRun: dryRun,
		})

		// get block hash and tx hash from args/flags
		blockHash := common.HexToHash(args[0])
		rawTxHash, _ := cmd.Flags().GetString("tx")
		txHash := common.HexToHash(rawTxHash)

		tx, err := d.DefenderDA(blockHash, txHash)
		if err != nil {
			logger.Error("Failed to defend data availability", "err", err)
			panic(err)
		}

		fmt.Println("Defended data availability with tx:", tx.Hash().Hex(), "gas used:", tx.Gas(), "gas price:", tx.GasPrice().Uint64())
	},
}
