package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	StorePath string `mapstructure:"storePath"`
	Celestia struct {
		DANodeRPC    string  `mapstructure:"daNodeRPC"`
		DANodeToken  string  `mapstructure:"daNodeToken"`
		Namespace    string  `mapstructure:"namespace"`
		ConsensusRPC string  `mapstructure:"consensusRPC"`
		GasPrice     float64 `mapstructure:"gasPrice"`
		Retries      int     `mapstructure:"retries"`
		RetryDelay   int     `mapstructure:"retryDelay"`
		Mnemonic     string  `mapstructure:"mnemonic"`
		ConsensusGRPC string `mapstructure:"consensusGRPC"`
		ConsensusTLS  bool   `mapstructure:"consensusTLS"`
		Network       string `mapstructure:"network"`
	} `mapstructure:"celestia"`
	Ethereum struct {
		HTTPEndpoint            string `mapstructure:"httpEndpoint"`
		WSEndpoint              string `mapstructure:"wsEndpoint"`
		CanonicalStateChain     string `mapstructure:"canonicalStateChain"`
		DaOracle                string `mapstructure:"daOracle"`
		GasPriceIncreasePercent int    `mapstructure:"gasPriceIncreasePercent"`
		Challenge               string `mapstructure:"challenge"`
		ChainOracle             string `mapstructure:"chainOracle"`
		BlobstreamX             string `mapstructure:"blobstreamX"`
		BlockTime               int    `mapstructure:"blockTime"`
		Timeout                 int    `mapstructure:"timeout"`
	} `mapstructure:"ethereum"`
	LightLink struct {
		Endpoint            string `mapstructure:"endpoint"`
		Delay               int    `mapstructure:"delay"`
		L2ToL1MessagePasser string `mapstructure:"l2ToL1MessagePasser"`
	} `mapstructure:"lightlink"`
	Rollup struct {
		L1PollDelay int    `mapstructure:"l1pollDelay"`
		L2PollDelay int    `mapstructure:"l2pollDelay"`
		BundleSize  uint64 `mapstructure:"bundleSize"`
		BundleCount uint64 `mapstructure:"bundleCount"`
		Store       bool   `mapstructure:"store"`
	} `mapstructure:"rollup"`
	Defender struct {
		WorkerDelay int `mapstructure:"workerDelay"`
	} `mapstructure:"defender"`
	// Not typically set in config file.
	DryRun bool `mapstructure:"dryRun,omitempty"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(viper.GetString("config-path"))
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
