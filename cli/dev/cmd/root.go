package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// variables for the root command
	cfgPath   string
	logLevel  string
	logType   string
	logSource bool

	// root command
	RootCmd = &cobra.Command{
		Use:   "hbdev",
		Short: "Developer CLI for Hummingbird",
	}
)

// init function to set up the root command
func init() {
	RootCmd.PersistentFlags().StringVar(&cfgPath, "config-path", ".", "sets the config file path (default is .)")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "sets the log output level (default is info)")
	RootCmd.PersistentFlags().StringVar(&logType, "log-type", "console", "sets the log output type [console,json] (default is console)")
	RootCmd.PersistentFlags().BoolVar(&logSource, "log-source", false, "log output source file (default is false)")

	// bind flags to viper
	viper.BindPFlag("config-path", RootCmd.PersistentFlags().Lookup("config-path"))
	viper.BindPFlag("log-level", RootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("log-type", RootCmd.PersistentFlags().Lookup("log-type"))
	viper.BindPFlag("log-source", RootCmd.PersistentFlags().Lookup("log-source"))

	// bind version to viper
	viper.SetDefault("version", "0.0.1")
}
