package node

import (
	"crypto/ecdsa"
	"hummingbird/config"
	canonicalstatechain "hummingbird/node/contracts/CanonicalStateChain.sol"
	"hummingbird/node/ethereum"
	"log/slog"
	"math/big"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type Node struct {
	ethereum.Ethereum
	Celestia
	LightLink

	Store KVStore
}

// NewFromConfig creates a new node from the given config.
func NewFromConfig(cfg *config.Config, logger *slog.Logger, ethKey *ecdsa.PrivateKey) (*Node, error) {

	logger.Info("Starting LightLink Hummingbird ("+viper.GetString("version")+")",
		"Go Version", runtime.Version(),
		"Operating System", runtime.GOOS,
		"Architecture", runtime.GOARCH)

	// log config file path
	logger.Info("Using config file", "path", viper.ConfigFileUsed())

	eth, err := ethereum.NewClient(ethereum.ClientOpts{
		Endpoint:                   cfg.Ethereum.HTTPEndpoint,
		CanonicalStateChainAddress: common.HexToAddress(cfg.Ethereum.CanonicalStateChain),
		ChallengeAddress:           common.HexToAddress(cfg.Ethereum.Challenge),
		ChainOracleAddress:         common.HexToAddress(cfg.Ethereum.ChainOracle),
		BlobstreamXAddress:         common.HexToAddress(cfg.Ethereum.BlobstreamX),
		Signer:                     ethKey,
		Logger:                     logger.With("ctx", "ethereum-http"),
		DryRun:                     cfg.DryRun,
		GasPriceIncreasePercent:    big.NewInt(int64(cfg.Ethereum.GasPriceIncreasePercent)),
		BlockTime:                  cfg.Ethereum.BlockTime,
		Timeout:                    time.Duration(cfg.Ethereum.Timeout) * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	cel, err := NewCelestiaClient(CelestiaClientOpts{
		Endpoint:                cfg.Celestia.Endpoint,
		Token:                   cfg.Celestia.Token,
		TendermintRPC:           cfg.Celestia.TendermintRPC,
		Namespace:               cfg.Celestia.Namespace,
		Logger:                  logger.With("ctx", "celestia"),
		GasPrice:                cfg.Celestia.GasPrice,
		GasPriceIncreasePercent: big.NewInt(int64(cfg.Celestia.GasPriceIncreasePercent)),
		GasAPI:                  cfg.Celestia.GasAPI,
		Retries:                 cfg.Celestia.Retries,
		RetryDelay:              time.Duration(cfg.Celestia.RetryDelay) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	ll, err := NewLightLinkClient(&LightLinkClientOpts{
		Endpoint: cfg.LightLink.Endpoint,
		Delay:    time.Duration(cfg.LightLink.Delay) * time.Millisecond,
		Logger:   logger.With("ctx", "lightlink"),
	})
	if err != nil {
		return nil, err
	}

	var store *LDBStore

	if cfg.Rollup.Store {
		store, err = NewLDBStore(cfg.StorePath)
		if err != nil {
			return nil, err
		}
	}

	logger.Info("Rollup Node created!", "dryRun", cfg.DryRun)

	logger.Info("Ethereum private key address", "address", crypto.PubkeyToAddress(ethKey.PublicKey).Hex())

	return &Node{
		Ethereum:  eth,
		Celestia:  cel,
		LightLink: ll,

		Store: store,
	}, nil
}

// GetDAPointer gets the Celestia pointer for the given rollup block hash.
func (n *Node) GetDAPointer(hash common.Hash) ([]*CelestiaPointer, error) {

	// TODO FETCH FROM LOCAL STORE!

	// pointer is not found in local store so get rollup header
	header, err := n.Ethereum.GetRollupHeaderByHash(hash)
	if err != nil {
		return nil, err
	}

	// get pointer from header
	pointers := make([]*CelestiaPointer, 0)
	for i := 0; i < len(header.CelestiaPointers); i++ {

		pointers = append(pointers, &CelestiaPointer{
			Height:     header.CelestiaPointers[i].Height,
			ShareStart: header.CelestiaPointers[i].ShareStart.Uint64(),
			ShareLen:   uint64(header.CelestiaPointers[i].ShareLen),
		})
	}

	return pointers, nil
}

func (n *Node) FetchRollupBlock(rblock common.Hash) (*canonicalstatechain.CanonicalStateChainHeader, []*Bundle, error) {
	header, err := n.Ethereum.GetRollupHeaderByHash(rblock)
	if err != nil {
		return nil, nil, err
	}

	bundles := make([]*Bundle, 0)
	for i := 0; i < len(header.CelestiaPointers); i++ {
		pointer := &CelestiaPointer{
			Height:     header.CelestiaPointers[i].Height,
			ShareStart: header.CelestiaPointers[i].ShareStart.Uint64(),
			ShareLen:   uint64(header.CelestiaPointers[i].ShareLen),
		}

		shares, err := n.Celestia.GetSharesByNamespace(pointer)
		if err != nil {
			return nil, nil, err
		}

		bundle, err := NewBundleFromShares(shares)
		if err != nil {
			return nil, nil, err
		}

		bundles = append(bundles, bundle)
	}

	return &header, bundles, nil
}

// Returns true if the given ethKey is the publisher set in CanonicalStateChain
func (n *Node) IsPublisher(ethKey *ecdsa.PrivateKey) bool {
	if ethKey == nil {
		panic("eth key is nil")
	}

	p, err := n.Ethereum.GetPublisher()
	if err != nil {
		panic(err)
	}

	// Get address of public key
	addr := crypto.PubkeyToAddress(ethKey.PublicKey)

	return p == addr
}
