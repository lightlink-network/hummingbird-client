package node

import (
	"crypto/ecdsa"
	"hummingbird/config"
	"log/slog"
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
		Signer:                     ethKey,
	})
	if err != nil {
		return nil, err
	}
	logger.Debug("Ethereum client successfully initialized")

	cel, err := NewCelestiaClient(CelestiaClientOpts{
		Endpoint:      cfg.Celestia.Endpoint,
		Token:         cfg.Celestia.Token,
		GRPC:          cfg.Celestia.GRPC,
		TendermintRPC: cfg.Celestia.TendermintRPC,
		Namespace:     cfg.Celestia.Namespace,
	})
	if err != nil {
		return nil, err
	}
	logger.Debug("Celestia client successfully initialized")

	ll, err := NewLightLinkClient(&LightLinkClientOpts{
		Endpoint: cfg.LightLink.Endpoint,
		Delay:    time.Duration(cfg.LightLink.Delay) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	logger.Debug("LightLink client successfully initialized")

	store, err := NewLDBStore(cfg.StorePath)
	if err != nil {
		return nil, err
	}

	return &Node{
		Ethereum:  eth,
		Celestia:  cel,
		LightLink: ll,

		Store: store,
	}, nil
}
