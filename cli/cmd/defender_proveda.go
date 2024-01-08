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
)

func init() {
	DefenderProveDaCmd.PersistentFlags().String("tx", "", "celestia tx hash in which data was submitted")
	DefenderProveDaCmd.Flags().Bool("json", false, "output proof in json format")
}

var DefenderProveDaCmd = &cobra.Command{
	Use:   "prove-da",
	Short: "prove-da will prove a data availability batch",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		logger := DefaultLogger()
		ethKey := getEthKey()

		n, err := node.NewFromConfig(cfg, logger, ethKey)
		utils.NoErr(err)

		d := defender.NewDefender(n, &defender.Opts{
			Logger: logger.With("ctx", "Defender"),
		})

		rawTxHash, err := cmd.Flags().GetString("tx")
		if err != nil {
			logger.Error("Missing required tx hash from flag", "err", err)
			panic(err)
		}

		txHash := common.HexToHash(rawTxHash)
		proof, err := d.ProveDA(txHash)
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

		fmt.Println("Proof:")
		fmt.Println("	Nonce:", proof.Nonce)
		fmt.Println("	Tuple.Height:", proof.Tuple.Height)
		fmt.Println("	Tuple.DataRoot:", common.Hash(proof.Tuple.DataRoot).Hex())
		fmt.Println("	WrappedProof:", hexutil.Encode(wrappedProof))
	},
}
