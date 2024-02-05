package cmd

import (
	"encoding/json"
	"fmt"
	"hummingbird/config"
	"hummingbird/defender"
	"hummingbird/node"
	"hummingbird/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	DefenderProveDaCmd.Flags().Bool("json", false, "output proof in json format")
	DefenderProveDaCmd.Flags().Bool("verify", false, "verify the proof against the L1 rollup contract")
}

var DefenderProveDaCmd = &cobra.Command{
	Use:   "prove-da",
	Short: "prove-da will prove a data availability batch",
	Args:  cobra.MinimumNArgs(1),
	ArgAliases: []string{
		"block",
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := GetLogger(viper.GetString("log-type"))
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		blockHash := common.HexToHash(args[0])
		proof, err := d.GetDAProof(blockHash)
		if err != nil {
			logger.Error("Failed to prove data availability", "err", err)
			panic(err)
		}

		if useJson, _ := cmd.Flags().GetBool("json"); useJson {
			buf, err := json.MarshalIndent(proof, "", "  ")
			utils.NoErr(err)
			fmt.Println(string(buf))
			return
		}

		wrappedProof, err := rlp.EncodeToBytes(proof.WrappedProof)
		utils.NoErr(err)

		fmt.Println(" ")
		fmt.Println("Proof:")
		fmt.Println("	Nonce:", proof.Nonce)
		fmt.Println("	Tuple.Height:", proof.Tuple.Height)
		fmt.Println("	Tuple.DataRoot:", common.Hash(proof.Tuple.DataRoot).Hex())
		fmt.Println("	WrappedProof:", hexutil.Encode(wrappedProof))
		fmt.Println(" ")

		if verify, _ := cmd.Flags().GetBool("verify"); !verify {
			return
		}

		// Verify the proof against the L1 rollup contract.
		verified, err := n.Ethereum.DAVerify(proof)
		if err != nil {
			logger.Error("Failed to verify proof", "err", err)
			panic(err)
		}

		fmt.Println(" ")
		fmt.Println("Verified:", verified)

	},
}
