package main

import (
	"hummingbird/cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgPath   string
	logLevel  string
	logType   string
	logSource bool
)

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
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config-path", ".", "sets the config file path (default is .)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "sets the log output level (default is info)")
	rootCmd.PersistentFlags().StringVar(&logType, "log-type", "console", "sets the log output type [console,json] (default is console)")
	rootCmd.PersistentFlags().BoolVar(&logSource, "log-source", false, "log output source file (default is false)")
	// bind flags to viper
	viper.BindPFlag("config-path", rootCmd.PersistentFlags().Lookup("config-path"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("log-type", rootCmd.PersistentFlags().Lookup("log-type"))
	viper.BindPFlag("log-source", rootCmd.PersistentFlags().Lookup("log-source"))
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
