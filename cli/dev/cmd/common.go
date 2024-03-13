package cmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"hummingbird/cli/hb/cmd"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/utils"
	"log/slog"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func makeNode() (*node.Node, *slog.Logger, error) {
	cfg := config.Load()
	log := cmd.ConsoleLogger()
	ethKey := getEthKey()

	n, err := node.NewFromConfig(cfg, log, ethKey)
	if err != nil {
		return nil, nil, err
	}

	return n, log, nil
}

func getEthKey() *ecdsa.PrivateKey {
	key := os.Getenv("ETH_KEY")
	if key == "" {
		return nil
	}

	ethKey, err := crypto.ToECDSA(hexutil.MustDecode(key))
	if err != nil {
		return nil
	}

	return ethKey
}

// panicErr panics if err is not nil, with an optional prefix
func panicErr(err error, prefix ...string) {
	if err != nil {
		if len(prefix) > 0 {
			panic(fmt.Errorf("%s: %w", prefix[0], err))
		}
		panic(err)
	}
}

func printJSON(v interface{}) {
	output, err := utils.PrepareTidyJSON(v)
	panicErr(err, "json output preparation failed")

	buf, err := json.MarshalIndent(&output, "", " ")
	panicErr(err, "json output formatting failed")

	fmt.Println(string(buf))
}

func printPretty(v interface{}) {
	out := utils.MarshalText(v)
	fmt.Println(strings.TrimSpace(string(out)))
}
