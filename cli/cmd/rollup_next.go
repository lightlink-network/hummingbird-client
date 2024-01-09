package cmd

import (
	"fmt"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/node/contracts"
	"hummingbird/rollup"
	"hummingbird/utils"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	RollupNextCmd.Flags().Bool("dry", false, "dry run will not submit the rollup block to the L1 rollup contract, and will not upload real data to celestia")
}

var RollupNextCmd = &cobra.Command{
	Use:   "next",
	Short: "next will rollup the next batch of L2 blocks",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := ConsoleLogger()
		ethKey := getEthKey()

		// is dry run enabled?
		dryRun, _ := cmd.Flags().GetBool("dry")
		cfg.DryRun = dryRun

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		r := rollup.NewRollup(n, &rollup.Opts{
			PollDelay:             time.Duration(cfg.Rollup.PollDelay) * time.Millisecond,
			BundleSize:            cfg.Rollup.BundleSize,
			StoreCelestiaPointers: cfg.Rollup.StoreCelestiaPointers,
			StoreHeaders:          cfg.Rollup.StoreHeaders,
			Logger:                logger.With("ctx", "Rollup"),
			DryRun:                dryRun,
		})

		// If dry run is enabled, swap out celestia with a mock celestia client.
		if dryRun {
			logger.Warn("DryRun is enabled, using mock celestia client")
			r.Celestia = node.NewCelestiaMock(cfg.Celestia.Namespace)
		}

		logger.Info("Rolling up next batch of L2 blocks")
		b, err := r.CreateNextBlock()
		if err != nil {
			logger.Error("Failed to rollup next batch of L2 blocks", "err", err)
			panic(err)
		}

		hash, err := contracts.HashCanonicalStateChainHeader(b.CanonicalStateChainHeader)
		utils.NoErr(err)

		// Print out the rollup block.
		fmt.Println(" ")
		fmt.Println("Rollup Block:")
		fmt.Println("	Epoch:", b.Epoch)
		fmt.Println("	L2Height:", b.L2Height)
		fmt.Println("	PrevHash:", common.BytesToHash(b.PrevHash[:]).Hex())
		fmt.Println("	StateRoot:", common.BytesToHash(b.CanonicalStateChainHeader.StateRoot[:]).Hex())
		fmt.Println("	BlockRoot:", common.BytesToHash(b.CanonicalStateChainHeader.BlockRoot[:]).Hex())
		fmt.Println("	TxRoot:", common.BytesToHash(b.CanonicalStateChainHeader.TxRoot[:]).Hex())
		fmt.Println("	Hash:", hash.Hex())
		fmt.Println("	Bundle Size:", len(b.Bundle.Blocks))
		fmt.Println("	Celestia Height:", b.CelestiaHeight)
		fmt.Println("	Celestia Data Root:", common.BytesToHash(b.CelestiaDataRoot[:]).Hex())
		fmt.Println("	Celestia Tx Hash:", b.CelestiaPointer.TxHash.Hex())
		fmt.Println(" ")

		// If dry run is enabled, exit.
		// if dryRun {
		// 	logger.Warn("Dry run enabled, not submitting rollup block to L1 rollup contract")
		// 	return
		// }

		logger.Info(("Submitting rollup block to L1 rollup contract"))
		tx, err := r.SubmitBlock(b)
		if err != nil {
			logger.Error("Failed to submit rollup block to L1 rollup contract", "err", err)
			panic(err)
		}

		logger.Info("Rollup block submitted to L1 rollup contract", "tx_hash", tx.Hash().Hex())

	},
}
