package config

import "github.com/spf13/viper"

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
		DaOracle            string `mapstructure:"daOracle"`
	} `mapstructure:"ethereum"`
	LightLink struct {
		Endpoint string `mapstructure:"endpoint"`
		Delay    int    `mapstructure:"delay"`
	} `mapstructure:"lightlink"`
	Rollup struct {
		L1PollDelay           int    `mapstructure:"l1pollDelay"`
		L2PollDelay           int    `mapstructure:"l2pollDelay"`
		BundleSize            uint64 `mapstructure:"bundleSize"`
		StoreCelestiaPointers bool   `mapstructure:"storeCelestiaPointers"`
		StoreHeaders          bool   `mapstructure:"storeHeaders"`
	} `mapstructure:"rollup"`

	// Not typically set in config file.
	DryRun bool `mapstructure:"dryRun,omitempty"`
}

func Load() *Config {
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
