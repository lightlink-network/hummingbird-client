package cmd

import (
	"fmt"
	"hummingbird/node"
	"hummingbird/node/contracts"
	canonicalstatechain "hummingbird/node/contracts/CanonicalStateChain.sol"
	"hummingbird/utils"
	"strconv"
	"strings"

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

		// 0. parse flags
		showHeader, _ := cmd.Flags().GetBool("header")
		showBundle, _ := cmd.Flags().GetBool("bundle")
		showStats, _ := cmd.Flags().GetBool("stats")
		showShares, _ := cmd.Flags().GetBool("shares")
		showTxns, _ := cmd.Flags().GetBool("txns")

		// 1. get the rollup header
		var hash common.Hash
		var header canonicalstatechain.CanonicalStateChainHeader

		// - If the first argument is a hash, get the header by hash
		if strings.HasPrefix(args[0], "0x") {
			hash = common.HexToHash(args[0])
			panicErr(err, "invalid block hash")
			log.Info("Fetching Rollup Block", "hash", hash)
			header, err = n.Ethereum.GetRollupHeaderByHash(hash)
			panicErr(err, "failed to get rollup header")
			hash, _ = contracts.HashCanonicalStateChainHeader(&header)

			// - Otherwise, get the header by number
		} else {
			index, err := strconv.ParseUint(args[0], 10, 64)
			panicErr(err, "invalid block number")
			log.Info("Fetching Rollup Block", "index", hash)
			header, err = n.Ethereum.GetRollupHeader(index)
			panicErr(err, "failed to get rollup header")
			hash, _ = contracts.HashCanonicalStateChainHeader(&header)
		}

		// 2. download the rollup bundles shares
		log.Debug("Downloading shares...", "source", "Celestia")
		shares, err := n.Celestia.GetSharesByNamespace(&node.CelestiaPointer{
			Height:     header.CelestiaHeight,
			ShareStart: header.CelestiaShareStart,
			ShareLen:   header.CelestiaShareLen,
		})
		panicErr(err, "failed to get shares")
		log.Debug("✔️  Got Shares", "count", len(shares))

		// 3. decode the rollup bundle
		log.Debug("Decoding bundle...")
		data := utils.ExtractDataFromShares(shares)
		b := &node.Bundle{}
		err = b.DecodeRLP(data)
		panicErr(err, "failed to decode bundle")
		log.Debug("✔️ Got Bundle", "blocks", len(b.Blocks))

		// 4. print the rollup block
		if showHeader {
			fmt.Println("\n---- Rollup Block Header ----")
			fmt.Printf("Block %s\n", hash.Hex())
			printPretty(&header)
		}

		if showBundle {
			fmt.Println("\n---- Rollup Block Bundle ----")
			for i, block := range b.Blocks {
				fmt.Printf("Block #%4d, hash: %s, number: %d, txns: %d\n", i, hashWithoutExtraData(block).Hex(), block.NumberU64(), len(block.Transactions()))

			}
		}

		if showTxns {
			fmt.Println("\n---- Rollup Block Transactions ----")
			for i, block := range b.Blocks {
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
			for i, share := range shares {
				fmt.Printf("Share #%4d: 0x%x\n", i, share.ToBytes())
			}
		}

		if showStats {
			txs := []*types.Transaction{}
			for _, block := range b.Blocks {
				txs = append(txs, block.Transactions()...)
			}

			fmt.Println("\n---- Rollup Block Stats ----")
			fmt.Printf("Bundle Size: %d bytes\n", len(data))
			fmt.Printf("Shares Count: %d\n", len(shares))
			fmt.Printf("Block Count: %d\n", len(b.Blocks))
			fmt.Printf("Avg. Block Size: %d bytes\n", len(data)/len(b.Blocks))
			fmt.Printf("Tx Count: %d\n", len(txs))
			fmt.Printf("Avg. Tx Count: %d\n", len(txs)/len(b.Blocks))
		}
	},
}

func hashWithoutExtraData(block *types.Block) common.Hash {
	header := block.Header()
	header.Extra = common.Hex2Bytes("0x")
	return header.Hash()
}
