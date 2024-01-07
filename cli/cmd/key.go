package cmd

import "github.com/spf13/cobra"

func inti() {
	keyCmd.PersistentFlags().String("pass", "", "passphrase to encrypt/decrypt private key")
}

var keyCmd = &cobra.Command{
	Use: "key",
}
