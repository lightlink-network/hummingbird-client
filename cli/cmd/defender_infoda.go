package cmd

import (
	"encoding/json"
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	DefenderInfoDaCmd.Flags().Bool("json", false, "output proof in json format")
}

var DefenderInfoDaCmd = &cobra.Command{
	Use:        "info-da",
	Short:      "info-da will provide info on an existing challenge",
	ArgAliases: []string{"block"},
	Args:       cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		blockHash := common.HexToHash(args[0])

		info, err := d.InfoDA(blockHash)
		if err != nil {
			logger.Error("Failed to get data availability info", "err", err)
			panic(err)
		}

		if useJson, _ := cmd.Flags().GetBool("json"); useJson {
			buf, err := json.MarshalIndent(info, "", "  ")
			utils.NoErr(err)
			fmt.Println(string(buf))
			return
		}

		fmt.Println("Data Availability Info")
		fmt.Println(" ")
		fmt.Println(utils.MarshalText(&info))
		fmt.Println(" ")
	},
}
