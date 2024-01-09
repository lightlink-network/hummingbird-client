package node

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"hummingbird/node/contracts"
	"hummingbird/utils"
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
	DAVerify(*CelestiaProof) (bool, error)                                                  // Check if the data availability layer is verified.
	Wait(txHash common.Hash) (*types.Receipt, error)                                        // Wait for the given transaction to be mined.
}

type EthereumClient struct {
	signer              *ecdsa.PrivateKey
	client              *ethclient.Client
	chainId             *big.Int
	canonicalStateChain *contracts.CanonicalStateChainContract
	daOracle            *contracts.DAOracleContract
	logger              *slog.Logger
	opts                *EthereumClientOpts
}

type EthereumClientOpts struct {
	Signer                     *ecdsa.PrivateKey
	Endpoint                   string
	CanonicalStateChainAddress common.Address
	DAOracleAddress            common.Address
	Logger                     *slog.Logger
	DryRun                     bool
}

// NewEthereumRPC returns a new EthereumRPC client.
func NewEthereumRPC(opts EthereumClientOpts) (*EthereumClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	client, err := ethclient.Dial(opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	canonicalStateChain, err := contracts.NewCanonicalStateChainContract(opts.CanonicalStateChainAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CanonicalStateChain: %w", err)
	}

	daOracle, err := contracts.NewDAOracleContract(opts.DAOracleAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAOracle: %w", err)
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
	if ok, _ := utils.IsContract(client, opts.DAOracleAddress); !ok {
		opts.Logger.Warn("contract not found for DAOracle at given Address", "address", opts.DAOracleAddress.Hex(), "endpoint", opts.Endpoint)
	}

	return &EthereumClient{
		signer:              opts.Signer,
		client:              client,
		chainId:             chainId,
		canonicalStateChain: canonicalStateChain,
		daOracle:            daOracle,
		logger:              opts.Logger,
		opts:                &opts,
	}, nil
}

func (e *EthereumClient) transactor() (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(e.signer, e.chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	opts.GasPrice = gasPrice

	// If dry run is enabled, don't send the transaction.
	if e.opts.DryRun {
		e.logger.Warn("DryRun is enabled, not sending transaction")
		opts.NoSend = true
	}

	return opts, nil
}

// GetRollupHead returns the latest rollup block header.
func (e *EthereumClient) GetRollupHead() (contracts.CanonicalStateChainHeader, error) {
	return e.canonicalStateChain.GetHead(nil)
}

// PushRollupHead pushes a new rollup block header.
func (e *EthereumClient) PushRollupHead(header *contracts.CanonicalStateChainHeader) (*types.Transaction, error) {

	transactor, err := e.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
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

	// 1. try to get the the tx, see if it is pending
	_, isPending, err := e.client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// 2. if it is pending, wait for it to be mined
	if isPending {
		time.Sleep(1 * time.Second)
		return e.Wait(txHash)
	}

	// 3. otherwise, if it is not pending, get the receipt
	return e.client.TransactionReceipt(context.Background(), txHash)
}

func (e *EthereumClient) DAVerify(proof *CelestiaProof) (bool, error) {
	return e.daOracle.VerifyAttestation(nil, proof.Nonce, *proof.Tuple, *proof.WrappedProof)
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

func (e *ethereumMock) DAVerify(proof *CelestiaProof) (bool, error) {
	return true, nil
}

func (e *ethereumMock) SimulateHeight(height uint64) {
	e.height = height
}
