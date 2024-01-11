package node

import (
	"crypto/ecdsa"
	"hummingbird/config"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Node struct {
	Ethereum
	Celestia
	LightLink

	Store KVStore
}

// NewFromConfig creates a new node from the given config.
func NewFromConfig(cfg *config.Config, logger *slog.Logger, ethKey *ecdsa.PrivateKey) (*Node, error) {

	eth, err := NewEthereumRPC(EthereumClientOpts{
		Endpoint:                   cfg.Ethereum.Endpoint,
		CanonicalStateChainAddress: common.HexToAddress(cfg.Ethereum.CanonicalStateChain),
		DAOracleAddress:            common.HexToAddress(cfg.Ethereum.DaOracle),
		ChallengeAddress:           common.HexToAddress(cfg.Ethereum.Challenge),
		Signer:                     ethKey,
		Logger:                     logger.With("ctx", "ethereum"),
		DryRun:                     cfg.DryRun,
		GasPriceIncreasePercent:    big.NewInt(int64(cfg.Ethereum.GasPriceIncreasePercent)),
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

	store, err := NewLDBStore(cfg.StorePath)
	if err != nil {
		return nil, err
	}

	logger.Info("Rollup Node created!", "dryRun", cfg.DryRun)
	return &Node{
		Ethereum:  eth,
		Celestia:  cel,
		LightLink: ll,

		Store: store,
	}, nil
}
