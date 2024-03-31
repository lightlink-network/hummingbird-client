package cmd

import (
	"encoding/json"
	"fmt"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/rollup"
	"hummingbird/utils"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RollupInfoCmd.Flags().Bool("json", false, "output info in json format")
	RollupInfoCmd.Flags().String("hash", "", "block hash to get info for")
	RollupInfoCmd.Flags().Uint64("num", 0, "block number to get info for")
	RollupInfoCmd.Flags().Bool("bundle", false, "get bundle info. Requires --hash or --num to be set")
}

var RollupInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "info will print information about the current rollup state",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()
		useJson, _ := cmd.Flags().GetBool("json")

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		r := rollup.NewRollup(n, &rollup.Opts{
			L1PollDelay:           time.Duration(cfg.Rollup.L1PollDelay) * time.Millisecond,
			L2PollDelay:           time.Duration(cfg.Rollup.L2PollDelay) * time.Millisecond,
			BundleSize:            cfg.Rollup.BundleSize,
			StoreCelestiaPointers: cfg.Rollup.StoreCelestiaPointers,
			StoreHeaders:          cfg.Rollup.StoreHeaders,
			Logger:                logger.With("ctx", "Rollup"),
		})

		var useHash bool
		var useNum bool
		var num uint64

		// is a hash specified?
		hash, err := cmd.Flags().GetString("hash")
		useHash = err == nil && hash != ""
		// is a number specified?
		if cmd.Flags().Changed("num") {
			num, err = cmd.Flags().GetUint64("num")
			useNum = err == nil
		}

		var blockHash common.Hash
		if useNum {
			h, err := r.Ethereum.GetRollupHeader(num)
			utils.NoErr(err)
			blockHash, err = r.Ethereum.HashHeader(&h)
			utils.NoErr(err)
		}
		if useHash {
			blockHash = common.HexToHash(hash)
		}

		// if a hash or number is specified, get info for the block with that hash
		if useHash || useNum {
			info, err := r.GetBlockInfo(blockHash)
			utils.NoErr(err)
			printInfo(info, useJson)

			// if showBundle flag is set, get showBundle info
			if showBundle, _ := cmd.Flags().GetBool("bundle"); showBundle {
				for _, p := range info.CanonicalStateChainHeader.CelestiaPointers {
					s, err := r.Celestia.GetShares(&node.CelestiaPointer{
						Height:     p.Height,
						ShareStart: p.ShareStart.Uint64(),
						ShareLen:   uint64(p.ShareLen),
					})
					utils.NoErr(err)

					bundle, err := node.NewBundleFromShares(s)
					utils.NoErr(err)

					printBundle(bundle)
				}
			}
			return
		}

		// warn user if bundle flag is set but no block is specified
		if bundle, _ := cmd.Flags().GetBool("bundle"); bundle {
			logger.Warn("Bundle flag is set but no block specified. No Bundle info will be shown")
		}

		// otherwise get info for the chain
		info, err := r.GetInfo()
		utils.NoErr(err)
		printInfo(info, useJson)
	},
}

func printInfo(info any, useJson bool) {
	if useJson {
		buf, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(buf))
		return
	}

	// otherwise print as pretty text
	fmt.Println(" ")
	fmt.Println(utils.MarshalText(info))
	fmt.Println(" ")
}

func printBundle(bundle *node.Bundle) {
	fmt.Println(" ")
	fmt.Println("Bundle:")
	fmt.Println(" Blocks:", len(bundle.Blocks))
	for _, b := range bundle.Blocks {
		fmt.Println("  →", "Index:", b.Number(), "Hash:", utils.HashWithoutExtraData(b).Hex())
	}
	fmt.Println(" ")
}
