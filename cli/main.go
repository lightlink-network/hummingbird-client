package main

import (
	"hummingbird/cli/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hummingbird",
	Short: "Hummingbird is LightLinks rollup node",
}

func main() {
	rootCmd.AddCommand(cmd.RollupCmd)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
