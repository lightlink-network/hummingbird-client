package node

import (
	"context"
	"crypto/ecdsa"
	"hummingbird/node/contracts"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Ethereum is an ethereum client.
// It Provides access to the Ethereum Network and methods for
// interacting with important contracts on the network including:
// - CanonicalStateChain.sol With with methods for getting and pushing
// rollup block headers.
type Ethereum interface {
	GetRollupHeight() (uint64, error)                                                       // Get the current rollup block height.
	GetHeight() (uint64, error)                                                             // Get the current block height of the Ethereum network.
	GetRollupHead() (contracts.CanonicalStateChainHeader, error)                            // Get the latest rollup block header in the CanonicalStateChain.sol contract.
	PushRollupHead(header *contracts.CanonicalStateChainHeader) (*types.Transaction, error) // Push a new rollup block header to the CanonicalStateChain.sol contract.
	GetRollupHeader(index uint64) (contracts.CanonicalStateChainHeader, error)              // Get the rollup block header at the given index from the CanonicalStateChain.sol contract.
	GetRollupHeaderByHash(hash common.Hash) (contracts.CanonicalStateChainHeader, error)    // Get the rollup block header with the given hash from the CanonicalStateChain.sol contract.
	Wait(txHash common.Hash) (*types.Receipt, error)                                        // Wait for the given transaction to be mined.
}

type EthereumClient struct {
	signer              *ecdsa.PrivateKey
	client              *ethclient.Client
	chainId             *big.Int
	canonicalStateChain *contracts.CanonicalStateChainContract
	logger              *slog.Logger
}

type EthereumClientOpts struct {
	Signer                     *ecdsa.PrivateKey
	Endpoint                   string
	CanonicalStateChainAddress common.Address
	Logger                     *slog.Logger
}

// NewEthereumRPC returns a new EthereumRPC client.
func NewEthereumRPC(opts EthereumClientOpts) (*EthereumClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	log := opts.Logger.With("func", "NewEthereumRPC")

	client, err := ethclient.Dial(opts.Endpoint)
	if err != nil {
		log.Error("Failed to connect to Ethereum", "error", err)
		return nil, err
	}

	canonicalStateChain, err := contracts.NewCanonicalStateChainContract(opts.CanonicalStateChainAddress, client)
	if err != nil {
		log.Error("Failed to connect to CanonicalStateChain", "error", err)
		return nil, err
	}

	chainId, err := client.ChainID(context.TODO())
	if err != nil {
		log.Error("Failed to get chainId", "error", err)
		return nil, err
	}

	log.Info("Connected to Ethereum", "chainId", chainId)

	return &EthereumClient{
		signer:              opts.Signer,
		client:              client,
		chainId:             chainId,
		canonicalStateChain: canonicalStateChain,
		logger:              opts.Logger,
	}, nil
}

func (e *EthereumClient) transactor() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(e.signer, e.chainId)
}

// GetRollupHead returns the latest rollup block header.
func (e *EthereumClient) GetRollupHead() (contracts.CanonicalStateChainHeader, error) {
	return e.canonicalStateChain.GetHead(nil)
}

// PushRollupHead pushes a new rollup block header.
func (e *EthereumClient) PushRollupHead(header *contracts.CanonicalStateChainHeader) (*types.Transaction, error) {
	log := e.logger.With("func", "PushRollupHead")

	transactor, err := e.transactor()
	if err != nil {
		log.Error("Failed to create transactor", "error", err)
		return nil, err
	}

	return e.canonicalStateChain.PushBlock(transactor, *header)
}

// GetRollupHeader returns the rollup block header at the given index.
func (e *EthereumClient) GetRollupHeader(index uint64) (contracts.CanonicalStateChainHeader, error) {
	return e.canonicalStateChain.GetBlock(nil, big.NewInt(int64(index)))
}

// GetRollupHeaderByHash returns the rollup block header with the given hash.
func (e *EthereumClient) GetRollupHeaderByHash(hash common.Hash) (contracts.CanonicalStateChainHeader, error) {
	return e.canonicalStateChain.Headers(nil, hash)
}

// GetRollupHeight returns the current rollup block height.
func (e *EthereumClient) GetRollupHeight() (uint64, error) {
	h, err := e.canonicalStateChain.ChainHead(nil)
	if err != nil {
		return 0, err
	}

	return h.Uint64(), nil
}

func (e *EthereumClient) GetHeight() (uint64, error) {
	return e.client.BlockNumber(context.Background())
}

func (e *EthereumClient) Wait(txHash common.Hash) (*types.Receipt, error) {
	log := e.logger.With("func", "Wait")

	// 1. try to get the the tx, see if it is pending
	_, isPending, err := e.client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		log.Error("Failed to get transaction", "error", err)
		return nil, err
	}

	// 2. if it is pending, wait for it to be mined
	if isPending {
		time.Sleep(1 * time.Second)
		return e.Wait(txHash)
	}

	// 3. otherwise, if it is not pending, get the receipt
	return e.client.TransactionReceipt(context.Background(), txHash)
}

// MOCK CLIENT FOR TESTING

type ethereumMock struct {
	rollupHeaders map[common.Hash]contracts.CanonicalStateChainHeader
	indexToHash   map[uint64]common.Hash
	head          int64
	height        uint64
}

// NewEthereumMock returns a new EthereumMock client. It is used for testing.
func NewEthereumMock(genisis *contracts.CanonicalStateChainHeader) *ethereumMock {

	e := &ethereumMock{
		rollupHeaders: make(map[common.Hash]contracts.CanonicalStateChainHeader),
		indexToHash:   make(map[uint64]common.Hash),
		head:          -1,
	}

	e.PushRollupHead(genisis)
	return e
}

func (e *ethereumMock) Wait(txHash common.Hash) (*types.Receipt, error) {
	return types.NewReceipt(txHash[:], false, 21000), nil
}

// GetRollupHead returns the latest rollup block header.
func (e *ethereumMock) GetRollupHead() (contracts.CanonicalStateChainHeader, error) {
	return e.GetRollupHeader(uint64(e.head))
}

// PushRollupHead pushes a new rollup block header.
func (e *ethereumMock) PushRollupHead(header *contracts.CanonicalStateChainHeader) (*types.Transaction, error) {
	index := e.head + 1

	hash, err := contracts.HashCanonicalStateChainHeader(header)
	if err != nil {
		return nil, err
	}

	e.height++
	e.head = index
	e.rollupHeaders[hash] = *header
	e.indexToHash[uint64(index)] = hash

	return types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), hash.Bytes()), nil
}

// GetRollupHeader returns the rollup block header at the given index.
func (e *ethereumMock) GetRollupHeader(index uint64) (contracts.CanonicalStateChainHeader, error) {
	hash, ok := e.indexToHash[index]
	if !ok {
		return contracts.CanonicalStateChainHeader{}, nil
	}

	return e.rollupHeaders[hash], nil
}

// GetRollupHeaderByHash returns the rollup block header with the given hash.
func (e *ethereumMock) GetRollupHeaderByHash(hash common.Hash) (contracts.CanonicalStateChainHeader, error) {
	return e.rollupHeaders[hash], nil
}

// GetRollupHeight returns the current rollup block height.
func (e *ethereumMock) GetRollupHeight() (uint64, error) {
	return uint64(e.head), nil
}

func (e *ethereumMock) GetHeight() (uint64, error) {
	return e.height, nil
}

func (e *ethereumMock) SimulateHeight(height uint64) {
	e.height = height
}
