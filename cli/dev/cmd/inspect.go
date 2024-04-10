package cmd

import (
	"fmt"
	"hummingbird/rollup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
)

func init() {
	// Add the inspect command to the root command
	RootCmd.AddCommand(InspectCmd)

	InspectCmd.Flags().Bool("header", true, "print the rollup header, default true")
	InspectCmd.Flags().Bool("bundle", false, "print the rollup bundle")
	InspectCmd.Flags().Bool("stats", false, "print the rollup stats")
	InspectCmd.Flags().Bool("shares", false, "print the rollup shares")
	InspectCmd.Flags().Bool("txns", false, "print the rollup transactions")
}

var InspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect will inspect a rollup block",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		n, log, err := makeNode()
		panicErr(err, "failed to create node")

		r := rollup.NewRollup(n, &rollup.Opts{
			Logger: log.With("ctx", "Rollup"),
		})

		// 0. parse flags
		showHeader, _ := cmd.Flags().GetBool("header")
		showBundle, _ := cmd.Flags().GetBool("bundle")
		showStats, _ := cmd.Flags().GetBool("stats")
		showShares, _ := cmd.Flags().GetBool("shares")
		showTxns, _ := cmd.Flags().GetBool("txns")

		// 1. get the rollup block
		hash := common.HexToHash(args[0])
		panicErr(err, "invalid block hash")
		log.Info("Fetching Rollup Block", "hash", hash)
		rblock, err := r.GetBlockByHash(hash)
		panicErr(err, "failed to get rollup block")

		// 4. print the rollup block
		if showHeader {
			fmt.Println("\n---- Rollup Block Header ----")
			fmt.Printf("Block %s\n", hash.Hex())
			printPretty(&rblock.CanonicalStateChainHeader)
		}

		if showBundle {
			fmt.Println("\n---- Rollup Block Bundle ----")
			for i, bundle := range rblock.Bundles {
				fmt.Printf("Bundle #%4d\n", i)
				for j, block := range bundle.Blocks {
					fmt.Printf("Block #%4d, hash: %s, number: %d, txns: %d\n", i+j, hashWithoutExtraData(block).Hex(), block.NumberU64(), len(block.Transactions()))
				}
			}
		}

		if showTxns {
			fmt.Println("\n---- Rollup Block Transactions ----")
			for i, block := range rblock.L2Blocks() {
				if len(block.Transactions()) > 0 {
					fmt.Printf("Block #%4d, hash: %s, number: %d\n", i, hashWithoutExtraData(block).Hex(), block.NumberU64())
					for j, tx := range block.Transactions() {
						fmt.Printf("  → Tx #%4d: (type: %d, size: %d) %s\n", j, tx.Type(), tx.Size(), tx.Hash().Hex())
					}
				}
			}
		}

		if showShares {
			fmt.Println("\n---- Rollup Block Shares ----")
			for i, bundle := range rblock.Bundles {
				fmt.Printf("Bundle #%4d Shares\n", i)
				ss, _ := bundle.Shares(n.Namespace())
				for j, share := range ss {
					fmt.Printf("Share #%4d: 0x%x\n", j, share.ToBytes())
				}
			}
		}

		if showStats {
			txs := []*types.Transaction{}
			for _, block := range rblock.L2Blocks() {
				txs = append(txs, block.Transactions()...)
			}

			fmt.Println("\n---- Rollup Block Stats ----")
			// fmt.Printf("Bundle Size: %d bytes\n", len(rblock))
			// fmt.Printf("Shares Count: %d\n", len(shares))
			fmt.Printf("Block Count: %d\n", len(rblock.L2Blocks()))
			// fmt.Printf("Avg. Block Size: %d bytes\n", len(data)/len(rblock.L2Blocks()))
			fmt.Printf("Tx Count: %d\n", len(txs))
			fmt.Printf("Avg. Tx Count: %d\n", len(txs)/len(rblock.L2Blocks()))
		}
	},
}

func hashWithoutExtraData(block *types.Block) common.Hash {
	header := block.Header()
	header.Extra = common.Hex2Bytes("0x")
	return header.Hash()
}
