package cmd

import (
	"crypto/ecdsa"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type Config struct {
	StorePath string `mapstructure:"storePath"`
	Celestia  struct {
		Token         string `mapstructure:"token"`
		Endpoint      string `mapstructure:"endpoint"`
		Namespace     string `mapstructure:"namespace"`
		TendermintRPC string `mapstructure:"tendermint_rpc"`
		GRPC          string `mapstructure:"grpc"`
	} `mapstructure:"celestia"`
	Ethereum struct {
		Endpoint            string `mapstructure:"endpoint"`
		CanonicalStateChain string `mapstructure:"canonicalStateChain"`
	} `mapstructure:"ethereum"`
	LightLink struct {
		Endpoint string `mapstructure:"endpoint"`
		Delay    int    `mapstructure:"delay"`
	} `mapstructure:"lightlink"`
	Rollup struct {
		PollDelay             int    `mapstructure:"pollDelay"`
		BundleSize            uint64 `mapstructure:"bundleSize"`
		StoreCelestiaPointers bool   `mapstructure:"storeCelestiaPointers"`
		StoreHeaders          bool   `mapstructure:"storeHeaders"`
	} `mapstructure:"rollup"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func DefaultLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(handler).With("app", "hummingbird")
	return logger
}

func getEthKey() *ecdsa.PrivateKey {
	key := os.Getenv("ETH_KEY")
	if key == "" {
		panic("env ETH_KEY not set")
	}

	ethKey, err := crypto.ToECDSA(hexutil.MustDecode(key))
	if err != nil {
		panic("Failed to decode ETH_KEY: " + err.Error())
	}

	return ethKey
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
