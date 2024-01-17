package cmd

import (
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DefenderStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start will start the defender node",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		err = d.Start()
		if err != nil {
			logger.Error("Defender.Start failed", "err", err, "retry_in", "5s")
		}

	},
}
