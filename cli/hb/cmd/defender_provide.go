package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	DefenderProvideCmd.Flags().String("type", "header", "type of data to provide (header, tx)")
	DefenderProvideCmd.Flags().Bool("skip-shares", false, "skip providing shares")
}

var DefenderProvideCmd = &cobra.Command{
	Use:   "provide",
	Short: "provide will download data from Celestia and provide it to Layer 1",
	Aliases: []string{
		"rblock",
		"hash",
	},
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()

		rblockHash := common.HexToHash(args[0])

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		// allow block hash or number
		targetHash := common.HexToHash(args[1])

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		// skip shares
		skipShares, _ := cmd.Flags().GetBool("skip-shares")
		if skipShares {
			logger.Info("Skipping providing shares")
		}

		// type
		t, _ := cmd.Flags().GetString("type")
		var tx *types.Transaction
		switch t {
		case "header":
			logger.Info("Providing L2 Header...")
			tx, err = d.ProvideL2Header(rblockHash, targetHash, skipShares)
			if err != nil {
				logger.Error("Defender.Provide header failed", "err", err)
				return
			}
		case "tx":
			logger.Info("Providing L2 Tx...")
			tx, err = d.ProvideL2Tx(rblockHash, targetHash, skipShares)
			if err != nil {
				logger.Error("Defender.Provide tx failed", "err", err)
				return
			}
		default:
			logger.Error("Invalid type", "type", t)
			return
		}

		fmt.Println(" ")
		fmt.Println("Tx Hash:", tx.Hash().Hex())
		fmt.Println("Provided L2 Data:", targetHash.Hex())
		fmt.Println("Included in Rollup Block:", rblockHash.Hex())
		fmt.Println(" ")
	},
}
