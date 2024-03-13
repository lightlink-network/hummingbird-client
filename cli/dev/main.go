package main

import "hummingbird/cli/dev/cmd"

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
