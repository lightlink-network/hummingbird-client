package cmd

import (
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var DefendHeaderCmd = &cobra.Command{
	Use:        "defend-header",
	Short:      "Defend L2 header",
	ArgAliases: []string{"rblock", "l2num"},
	Args:       cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := ConsoleLogger()
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		// get block hash and l2 num from args
		rblockHash := common.HexToHash(args[0])
		l2num, _ := new(big.Int).SetString(args[1], 10)
		logger.Info("Defending L2 header", "rblock", rblockHash.Hex(), "l2num", l2num.String())

		tx, err := d.DefendL2Header(rblockHash, l2num)
		utils.NoErr(err)

		logger.Info("Defended L2 header", "tx", tx.Hash().Hex())
	},
}
