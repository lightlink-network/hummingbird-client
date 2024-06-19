package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"hummingbird/utils"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	blobstreamXContract "hummingbird/node/contracts/BlobstreamX.sol"
	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
	chainOracleContract "hummingbird/node/contracts/ChainOracle.sol"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
)

// Ethereum is an ethereum client.
// It Provides access to the Ethereum Network and methods for
// interacting with important contracts on the network including:
// - CanonicalStateChain.sol With with methods for getting and pushing rollup block headers.
// - Challenge.sol With methods for challenging data availability etc
// - ChainOracle.sol With methods for providing shares and headers
// - BlobstreamX.sol With methods for verifying data availability
type Ethereum interface {
	CanonicalStateChain
	Challenge
	ChainOracle
	BlobstreamX
}

type Client struct {
	signer              *ecdsa.PrivateKey
	client              *ethclient.Client
	chainId             *big.Int
	canonicalStateChain *canonicalStateChainContract.CanonicalStateChain
	challenge           *challengeContract.Challenge
	chainLoader         *chainOracleContract.ChainOracle
	blobstreamX         *blobstreamXContract.BlobstreamX
	logger              *slog.Logger
	opts                *ClientOpts
}

type ClientOpts struct {
	Signer                     *ecdsa.PrivateKey
	Endpoint                   string
	CanonicalStateChainAddress common.Address
	ChallengeAddress           common.Address
	ChainOracleAddress         common.Address
	BlobstreamXAddress         common.Address
	Logger                     *slog.Logger
	DryRun                     bool
	GasPriceIncreasePercent    *big.Int
	BlockTime                  int
	Timeout                    time.Duration
}

// NewEthereumRPC returns a new EthereumRPC client over HTTP.
func NewClient(opts ClientOpts) (*Client, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	client, err := ethclient.Dial(opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	canonicalStateChain, err := canonicalStateChainContract.NewCanonicalStateChain(opts.CanonicalStateChainAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CanonicalStateChain: %w", err)
	}

	challenge, err := challengeContract.NewChallenge(opts.ChallengeAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Challenge: %w", err)
	}

	chainLoader, err := chainOracleContract.NewChainOracle(opts.ChainOracleAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ChainOracle: %w", err)
	}

	blobstreamX, err := blobstreamXContract.NewBlobstreamX(opts.BlobstreamXAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to BlobstreamX: %w", err)
	}

	chainId, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get chainId: %w", err)
	}

	opts.Logger.Info("Connected to Ethereum", "chainId", chainId)

	// Warn user if the contracts are not found at the given addresses.
	if ok, _ := utils.IsContract(client, opts.CanonicalStateChainAddress); !ok {
		opts.Logger.Warn("contract not found for CanonicalStateChain at given Address", "address", opts.CanonicalStateChainAddress.Hex(), "endpoint", opts.Endpoint)
	}
	if ok, _ := utils.IsContract(client, opts.ChallengeAddress); !ok {
		opts.Logger.Warn("contract not found for Challenge at given Address", "address", opts.ChallengeAddress.Hex(), "endpoint", opts.Endpoint)
	}
	if ok, _ := utils.IsContract(client, opts.ChainOracleAddress); !ok {
		opts.Logger.Warn("contract not found for ChainOracle at given Address", "address", opts.ChainOracleAddress.Hex(), "endpoint", opts.Endpoint)
	}
	if ok, _ := utils.IsContract(client, opts.BlobstreamXAddress); !ok {
		opts.Logger.Warn("contract not found for BlobstreamX at given Address", "address", opts.BlobstreamXAddress.Hex(), "endpoint", opts.Endpoint)
	}

	return &Client{
		signer:              opts.Signer,
		client:              client,
		chainId:             chainId,
		canonicalStateChain: canonicalStateChain,
		challenge:           challenge,
		chainLoader:         chainLoader,
		blobstreamX:         blobstreamX,
		logger:              opts.Logger,
		opts:                &opts,
	}, nil
}

func (e *Client) transactor() (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(e.signer, e.chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	opts.GasPrice = gasPrice

	// If gas price increase percent is set, increase the gas price by the given percent.
	if e.opts.GasPriceIncreasePercent != nil && e.opts.GasPriceIncreasePercent.Cmp(big.NewInt(0)) > 0 {
		opts.GasPrice = gasPrice.Add(gasPrice, new(big.Int).Div(new(big.Int).Mul(gasPrice, e.opts.GasPriceIncreasePercent), big.NewInt(100)))
	}

	// If dry run is enabled, don't send the transaction.
	if e.opts.DryRun {
		e.logger.Warn("DryRun is enabled, not sending transaction")
		opts.NoSend = true
	}

	return opts, nil
}
