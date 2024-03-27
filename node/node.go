package node

import (
	"crypto/ecdsa"
	"hummingbird/config"
	canonicalstatechain "hummingbird/node/contracts/CanonicalStateChain.sol"
	"log/slog"
	"math/big"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type Node struct {
	Ethereum
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

	eth, err := NewEthereumRPC(EthereumHTTPClientOpts{
		Endpoint:                   cfg.Ethereum.HTTPEndpoint,
		CanonicalStateChainAddress: common.HexToAddress(cfg.Ethereum.CanonicalStateChain),
		ChallengeAddress:           common.HexToAddress(cfg.Ethereum.Challenge),
		ChainOracleAddress:         common.HexToAddress(cfg.Ethereum.ChainOracle),
		BlobstreamXAddress:         common.HexToAddress(cfg.Ethereum.BlobstreamX),
		Signer:                     ethKey,
		Logger:                     logger.With("ctx", "ethereum-http"),
		DryRun:                     cfg.DryRun,
		GasPriceIncreasePercent:    big.NewInt(int64(cfg.Ethereum.GasPriceIncreasePercent)),
	}, EthereumWSClientOpts{
		Endpoint:         cfg.Ethereum.WSEndpoint,
		ChallengeAddress: common.HexToAddress(cfg.Ethereum.Challenge),
		Logger:           logger.With("ctx", "ethereum-ws"),
	})
	if err != nil {
		return nil, err
	}

	cel, err := NewCelestiaClient(CelestiaClientOpts{
		Endpoint:      cfg.Celestia.Endpoint,
		Token:         cfg.Celestia.Token,
		GRPC:          cfg.Celestia.GRPC,
		TendermintRPC: cfg.Celestia.TendermintRPC,
		Namespace:     cfg.Celestia.Namespace,
		Logger:        logger.With("ctx", "celestia"),
		GasPrice:      cfg.Celestia.GasPrice,
		Retries:       cfg.Celestia.Retries,
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

	if cfg.Rollup.StoreHeaders || cfg.Rollup.StoreCelestiaPointers {
		store, err = NewLDBStore(cfg.StorePath)
		if err != nil {
			return nil, err
		}
	}

	logger.Info("Rollup Node created!", "dryRun", cfg.DryRun)
	return &Node{
		Ethereum:  eth,
		Celestia:  cel,
		LightLink: ll,

		Store: store,
	}, nil
}

// GetDAPointer gets the Celestia pointer for the given rollup block hash.
func (n *Node) GetDAPointer(hash common.Hash) (*CelestiaPointer, error) {
	pointer := &CelestiaPointer{}

	if config.Load().Rollup.StoreCelestiaPointers {
		pointer, err := n.Store.GetDAPointer(hash)
		// if err is not found, get pointer from header, any other error return
		if err != nil && err.Error() != "failed to get celestia pointer from store: leveldb: not found" {
			return nil, err
		}

		// if pointer is found, return it
		if pointer != nil {
			return pointer, nil
		}
	}

	// pointer is not found in local store so get rollup header
	header, err := n.GetRollupHeaderByHash(hash)
	if err != nil {
		return nil, err
	}

	// get pointer from header
	pointer = &CelestiaPointer{
		Height:     header.CelestiaHeight,
		ShareStart: header.CelestiaShareStart,
		ShareLen:   header.CelestiaShareLen,
	}

	return pointer, nil
}

func (n *Node) FetchRollupBlock(rblock common.Hash) (*canonicalstatechain.CanonicalStateChainHeader, *Bundle, error) {
	header, err := n.GetRollupHeaderByHash(rblock)
	if err != nil {
		return nil, nil, err
	}

	pointer := &CelestiaPointer{
		Height:     header.CelestiaHeight,
		ShareStart: header.CelestiaShareStart,
		ShareLen:   header.CelestiaShareLen,
	}

	shares, err := n.Celestia.GetSharesByNamespace(pointer)
	if err != nil {
		return nil, nil, err
	}

	bundle, err := NewBundleFromShares(shares)
	if err != nil {
		return nil, nil, err
	}

	return &header, bundle, nil
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
