package cmd

import (
	"fmt"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

var formatPointer string
var verifyPointer bool

func init() {
	// Add the pointer command to the root command
	RootCmd.AddCommand(PointerCmd)

	PointerCmd.Flags().StringVar(&formatPointer, "format", "pretty", "output format [json, pretty]")
	PointerCmd.Flags().BoolVar(&verifyPointer, "verify", false, "verify the data pointer")
}

var PointerCmd = &cobra.Command{
	Use:        "pointer",
	Short:      "pointer finds the Celestia data pointer for a given hash",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"tx-hash"},
	Run: func(cmd *cobra.Command, args []string) {
		n, log, err := makeNode()
		panicErr(err, "failed to create node")

		// 0. parse flags
		txHash, err := hexutil.Decode(args[0])
		panicErr(err, "invalid celestia transaction hash")

		// 1. get the data pointer
		log.Info("Fetching celestia pointer", "tx", common.Hash(txHash))
		dataPointer, err := n.Celestia.GetPointer(common.BytesToHash(txHash))
		panicErr(err, "failed to get data pointer")

		// 2. verify the data pointer
		if verifyPointer {
			log.Info("Verifying Data Pointer", "celestiaHeight", dataPointer.Height, "shareStart", dataPointer.ShareStart, "shareLength", dataPointer.ShareLen)
			shares, err := n.Celestia.GetSharesByPointer(dataPointer)
			panicErr(err, "failed to get shares")

			data := utils.ExtractDataFromShares(shares)
			b := &node.Bundle{}
			err = b.DecodeRLP(data)
			panicErr(err, "failed to decode bundle")

			log.Info("✔️ Data Pointer Verified", "bundleSize", b.Size())

			// set the share root
			dataPointer.ShareRoot = b.BlockRoot()
		}

		// 3. print the data pointer
		fmt.Println("  ")
		printPretty(dataPointer)
		fmt.Println("  ")
	},
}
