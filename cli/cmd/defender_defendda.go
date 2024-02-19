package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	DefenderDefendDaCmd.Flags().Bool("dry", false, "dry run will not submit the rollup block to the L1 rollup contract, and will not upload real data to celestia")
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

		if dryRun {
			logger.Warn("DryRun is enabled, using mock celestia client")
			celestiaMock := node.NewCelestiaMock(cfg.Celestia.Namespace)
			celestiaMock.SetFakeProof(true)
			n.Celestia = celestiaMock
		}

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		// get block hash and tx hash from args/flags
		blockHash := common.HexToHash(args[0])

		tx, err := d.DefendDA(blockHash)
		if err != nil {
			if strings.Contains(err.Error(), "no data commitment has been generated for the provided height") {
				logger.Error("Failed to defend data availability, please wait for Celestia validators to commit data root", "err", err)
				return
			}

			logger.Error("Failed to defend data availability", "err", err)
			return
		}

		fmt.Println("Defended data availability with tx:", tx.Hash().Hex(), "gas used:", tx.Gas(), "gas price:", tx.GasPrice().Uint64())
	},
}
