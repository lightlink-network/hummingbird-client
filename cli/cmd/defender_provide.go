package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DefenderProvideCmd = &cobra.Command{
	Use:   "provide",
	Short: "provide will download data from Celestia and provide them to Layer 1",
	Aliases: []string{
		"rblock",
		"l2block",
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
		var l2blockHash common.Hash
		if strings.HasPrefix(args[1], "0x") {
			l2blockHash = common.HexToHash(args[1])
		} else {
			logger.Info("Providing L2 Header by block number", "block", args[1])
			num, err := strconv.Atoi(args[1])
			utils.NoErr(err)
			b, err := n.LightLink.GetBlock(uint64(num))
			utils.NoErr(err)

			l2blockHash = b.Header().Hash()
			logger.Info("Providing L2 Header by block number", "block", args[1], "hash", l2blockHash.Hex())
		}

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		tx, err := d.ProvideL2Header(rblockHash, l2blockHash)
		if err != nil {
			logger.Error("Defender.Provide failed", "err", err)
		}

		fmt.Println(" ")
		fmt.Println("Tx Hash:", tx.Hash().Hex())
		fmt.Println("Provided L2 Header:", l2blockHash.Hex())
		fmt.Println("Included in Rollup Block:", rblockHash.Hex())
		fmt.Println(" ")
	},
}
