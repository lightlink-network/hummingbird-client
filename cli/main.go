package main

import (
	"hummingbird/cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "hb",
	Short: "Hummingbird is LightLinks rollup node. It can be used to create new rollup blocks, download state, create and respond to challenges, and more.",
}

var rollupCmd = &cobra.Command{
	Use:   "rollup",
	Short: "rollup is a command to interact with the LL rollups",
}

var defenderCmd = &cobra.Command{
	Use:   "defender",
	Short: "defender is a command to generate proofs and respond to challenges",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".", "config.json file path (default is .)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func main() {
	// add subcommands to defender
	defenderCmd.AddCommand(cmd.DefenderProveDaCmd)

	// add subcommands to rollup
	rollupCmd.AddCommand(cmd.RollupInfoCmd)
	rollupCmd.AddCommand(cmd.RollupNextCmd)
	rollupCmd.AddCommand(cmd.RollupStartCmd)

	// add all commands to root
	rootCmd.AddCommand(rollupCmd)
	rootCmd.AddCommand(defenderCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
