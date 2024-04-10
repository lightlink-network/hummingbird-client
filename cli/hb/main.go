package main

import (
	"fmt"
	"hummingbird/cli/hb/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version will be set at build time
var Version = "development"

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

var challengerCmd = &cobra.Command{
	Use:   "challenger",
	Short: "challenger is a command to create challenges",
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
	// bind version to viper
	viper.SetDefault("version", Version)
}

func main() {
	logo()

	// add subcommands to challenger
	challengerCmd.AddCommand(cmd.ChallengerChallengedaCmd)

	// add subcommands to defender
	defenderCmd.AddCommand(cmd.DefenderProveDaCmd)
	defenderCmd.AddCommand(cmd.DefenderDefendDaCmd)
	defenderCmd.AddCommand(cmd.DefenderInfoDaCmd)
	defenderCmd.AddCommand(cmd.DefenderStartCmd)
	defenderCmd.AddCommand(cmd.DefenderProvideCmd)
	defenderCmd.AddCommand(cmd.DefendHeaderCmd)

	// add subcommands to rollup
	rollupCmd.AddCommand(cmd.RollupInfoCmd)
	rollupCmd.AddCommand(cmd.RollupNextCmd)
	rollupCmd.AddCommand(cmd.RollupStartCmd)

	// add all commands to root
	rootCmd.AddCommand(rollupCmd)
	rootCmd.AddCommand(defenderCmd)
	rootCmd.AddCommand(challengerCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func logo() {
	logo := `
				 _    _                           _             _     _         _ 
				| |  | |                         (_)           | |   (_)       | |
				| |__| |_   _ _ __ ___  _ __ ___  _ _ __   __ _| |__  _ _ __ __| |
				|  __  | | | | '_   _ \| '_   _ \| | '_ \ / _  | '_ \| | '__/ _  |
				| |  | | |_| | | | | | | | | | | | | | | | (_| | |_) | | | | (_| |
				|_|  |_|\__,_|_| |_| |_|_| |_| |_|_|_| |_|\__, |_.__/|_|_|  \__,_|
				                                           __/ |
				                                          |___/`

	fmt.Println(logo + "\n")
}
