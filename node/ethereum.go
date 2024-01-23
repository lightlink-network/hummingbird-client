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
	"github.com/ethereum/go-ethereum/event"

	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
	daOracleContract "hummingbird/node/contracts/DAOracle.sol"
)

// Ethereum is an ethereum client.
// It Provides access to the Ethereum Network and methods for
// interacting with important contracts on the network including:
// - CanonicalStateChain.sol With with methods for getting and pushing
// rollup block headers.
type Ethereum interface {
	GetRollupHeight() (uint64, error)                                                                         // Get the current rollup block height.
	GetHeight() (uint64, error)                                                                               // Get the current block height of the Ethereum network.
	GetRollupHead() (canonicalStateChainContract.CanonicalStateChainHeader, error)                            // Get the latest rollup block header in the CanonicalStateChain.sol contract.
	PushRollupHead(header *canonicalStateChainContract.CanonicalStateChainHeader) (*types.Transaction, error) // Push a new rollup block header to the CanonicalStateChain.sol contract.
	GetRollupHeader(index uint64) (canonicalStateChainContract.CanonicalStateChainHeader, error)              // Get the rollup block header at the given index from the CanonicalStateChain.sol contract.
	GetRollupHeaderByHash(hash common.Hash) (canonicalStateChainContract.CanonicalStateChainHeader, error)    // Get the rollup block header with the given hash from the CanonicalStateChain.sol contract.
	Wait(txHash common.Hash) (*types.Receipt, error)                                                          // Wait for the given transaction to be mined.
	DAVerify(*CelestiaProof) (bool, error)
	// Check if the data availability layer is verified.
	// Challenges
	GetChallengeFee() (*big.Int, error)
	GetDataRootInclusionChallenge(block common.Hash) (contracts.ChallengeDaInfo, error)
	ChallengeDataRootInclusion(index uint64) (*types.Transaction, common.Hash, error)
	DefendDataRootInclusion(common.Hash, *CelestiaProof) (*types.Transaction, error)
	SettleDataRootInclusion(common.Hash) (*types.Transaction, error)
	WatchChallengesDA(c chan<- *challengeContract.ChallengeChallengeDAUpdate, startBlock uint64) (event.Subscription, error)
}

type EthereumClient struct {
	http EthereumHTTPClient
	ws   EthereumWSClient
}

type EthereumHTTPClient struct {
	signer              *ecdsa.PrivateKey
	client              *ethclient.Client
	chainId             *big.Int
	canonicalStateChain *canonicalStateChainContract.CanonicalStateChain
	daOracle            *daOracleContract.DAOracleContract
	challenge           *challengeContract.Challenge
	logger              *slog.Logger
	opts                *EthereumHTTPClientOpts
}

type EthereumHTTPClientOpts struct {
	Signer                     *ecdsa.PrivateKey
	Endpoint                   string
	CanonicalStateChainAddress common.Address
	DAOracleAddress            common.Address
	ChallengeAddress           common.Address
	Logger                     *slog.Logger
	DryRun                     bool
	GasPriceIncreasePercent    *big.Int
}

type EthereumWSClient struct {
	client    *ethclient.Client
	challenge *challengeContract.Challenge
	logger    *slog.Logger
	opts      *EthereumWSClientOpts
}

type EthereumWSClientOpts struct {
	Endpoint         string
	ChallengeAddress common.Address
	Logger           *slog.Logger
}

func NewEthereumRPC(httpOpts EthereumHTTPClientOpts, wsOpts EthereumWSClientOpts) (*EthereumClient, error) {
	http, err := NewEthereumHTTP(httpOpts)
	if err != nil {
		return nil, err
	}

	ws, err := NewEthereumWS(wsOpts)
	if err != nil {
		return nil, err
	}

	return &EthereumClient{
		http: *http,
		ws:   *ws,
	}, nil
}

// NewEthereumRPC returns a new EthereumRPC client over HTTP.
func NewEthereumHTTP(opts EthereumHTTPClientOpts) (*EthereumHTTPClient, error) {
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

	daOracle, err := daOracleContract.NewDAOracleContract(opts.DAOracleAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAOracle: %w", err)
	}

	challenge, err := challengeContract.NewChallenge(opts.ChallengeAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Challenge: %w", err)
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
	if ok, _ := utils.IsContract(client, opts.ChallengeAddress); !ok {
		opts.Logger.Warn("contract not found for Challenge at given Address", "address", opts.ChallengeAddress.Hex(), "endpoint", opts.Endpoint)
	}

	return &EthereumHTTPClient{
		signer:              opts.Signer,
		client:              client,
		chainId:             chainId,
		canonicalStateChain: canonicalStateChain,
		daOracle:            daOracle,
		challenge:           challenge,
		logger:              opts.Logger,
		opts:                &opts,
	}, nil
}

// NewEthereumWS returns a new EthereumWS client over WebSockets.
func NewEthereumWS(opts EthereumWSClientOpts) (*EthereumWSClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	client, err := ethclient.Dial(opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum WebSocket: %w", err)
	}

	challenge, err := challengeContract.NewChallenge(opts.ChallengeAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Challenge contract: %w", err)
	}

	return &EthereumWSClient{
		client:    client,
		challenge: challenge,
		logger:    opts.Logger,
		opts:      &opts,
	}, nil
}

func (e *EthereumClient) transactor() (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(e.http.signer, e.http.chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	gasPrice, err := e.http.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	opts.GasPrice = gasPrice

	// If gas price increase percent is set, increase the gas price by the given percent.
	if e.http.opts.GasPriceIncreasePercent != nil && e.http.opts.GasPriceIncreasePercent.Cmp(big.NewInt(0)) > 0 {
		opts.GasPrice = gasPrice.Add(gasPrice, new(big.Int).Div(new(big.Int).Mul(gasPrice, e.http.opts.GasPriceIncreasePercent), big.NewInt(100)))
	}

	// If dry run is enabled, don't send the transaction.
	if e.http.opts.DryRun {
		e.http.logger.Warn("DryRun is enabled, not sending transaction")
		opts.NoSend = true
	}

	return opts, nil
}

// GetRollupHead returns the latest rollup block header.
func (e *EthereumClient) GetRollupHead() (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return e.http.canonicalStateChain.GetHead(nil)
}

// PushRollupHead pushes a new rollup block header.
func (e *EthereumClient) PushRollupHead(header *canonicalStateChainContract.CanonicalStateChainHeader) (*types.Transaction, error) {

	transactor, err := e.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return e.http.canonicalStateChain.PushBlock(transactor, *header)
}

// GetRollupHeader returns the rollup block header at the given index.
func (e *EthereumClient) GetRollupHeader(index uint64) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return e.http.canonicalStateChain.GetBlock(nil, big.NewInt(int64(index)))
}

// GetRollupHeaderByHash returns the rollup block header with the given hash.
func (e *EthereumClient) GetRollupHeaderByHash(hash common.Hash) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return e.http.canonicalStateChain.Headers(nil, hash)
}

// GetRollupHeight returns the current rollup block height.
func (e *EthereumClient) GetRollupHeight() (uint64, error) {
	h, err := e.http.canonicalStateChain.ChainHead(nil)
	if err != nil {
		return 0, err
	}

	return h.Uint64(), nil
}

func (e *EthereumClient) GetHeight() (uint64, error) {
	return e.http.client.BlockNumber(context.Background())
}

func (e *EthereumClient) Wait(txHash common.Hash) (*types.Receipt, error) {

	// 1. try to get the the tx, see if it is pending
	_, isPending, err := e.http.client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// 2. if it is pending, wait for it to be mined
	if isPending {
		time.Sleep(1 * time.Second)
		return e.Wait(txHash)
	}

	// 3. otherwise, if it is not pending, get the receipt
	return e.http.client.TransactionReceipt(context.Background(), txHash)
}

func (e *EthereumClient) DAVerify(proof *CelestiaProof) (bool, error) {
	// convert proof to daOracle format
	wrappedProof := daOracleContract.BinaryMerkleProof{
		SideNodes: proof.WrappedProof.SideNodes,
		Key:       proof.WrappedProof.Key,
		NumLeaves: proof.WrappedProof.NumLeaves,
	}
	return e.http.daOracle.VerifyAttestation(nil, proof.Nonce, *proof.Tuple, wrappedProof)
}

func (e *EthereumClient) GetChallengeFee() (*big.Int, error) {
	return e.http.challenge.ChallengeFee(nil)
}

func (e *EthereumClient) ChallengeDataRootInclusion(index uint64) (*types.Transaction, common.Hash, error) {
	transactor, err := e.transactor()
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to create transactor: %w", err)
	}

	// set transactions fee
	fee, err := e.GetChallengeFee()
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to get challenge fee: %w", err)
	}
	transactor.Value = fee

	// get index hash
	blockHash, err := e.http.canonicalStateChain.Chain(nil, big.NewInt(int64(index)))
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to get hash for block %d: %w", index, err)
	}

	tx, err := e.http.challenge.ChallengeDataRootInclusion(transactor, big.NewInt(int64(index)))
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to challenge data root inclusion: %w", err)
	}

	return tx, blockHash, nil
}

func (e *EthereumClient) DefendDataRootInclusion(blockHash common.Hash, proof *CelestiaProof) (*types.Transaction, error) {
	transactor, err := e.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := e.http.challenge.DefendDataRootInclusion(transactor, blockHash, challengeContract.ChallengeDataAvailabilityChallengeDAProof{
		RootNonce: proof.Nonce,
		DataRootTuple: challengeContract.DataRootTuple{
			Height:   proof.Tuple.Height,
			DataRoot: proof.Tuple.DataRoot,
		},
		Proof: *proof.WrappedProof,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to defend data root inclusion: %w", err)
	}

	return tx, nil
}

func (e *EthereumClient) SettleDataRootInclusion(blockHash common.Hash) (*types.Transaction, error) {
	transactor, err := e.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := e.http.challenge.SettleDataRootInclusion(transactor, blockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to settle data root inclusion: %w", err)
	}

	return tx, nil
}

func (e *EthereumClient) GetDataRootInclusionChallenge(blockHash common.Hash) (contracts.ChallengeDaInfo, error) {
	res, err := e.http.challenge.DaChallenges(nil, blockHash)
	if err != nil {
		return contracts.ChallengeDaInfo{}, fmt.Errorf("failed to get data root inclusion challenge: %w", err)
	}

	return contracts.ChallengeDaInfo{
		BlockIndex: res.BlockIndex,
		Challenger: res.Challenger.Hex(),
		Expiry:     res.Expiry,
		Status:     res.Status,
	}, nil
}

func (e *EthereumClient) WatchChallengesDA(c chan<- *challengeContract.ChallengeChallengeDAUpdate, startBlock uint64) (event.Subscription, error) {
	opts := &bind.WatchOpts{}
	blockHash := make([][32]byte, 0)
	blockIndex := make([]*big.Int, 0)
	statuses := make([]uint8, 0)

	// Create a new bind.WatchOpts that starts from the next block
	if startBlock > 0 {
		opts.Start = &startBlock
	}

	return e.ws.challenge.WatchChallengeDAUpdate(opts, c, blockHash, blockIndex, statuses)
}

// MOCK CLIENT FOR TESTING

type ethereumMock struct {
	rollupHeaders map[common.Hash]canonicalStateChainContract.CanonicalStateChainHeader
	indexToHash   map[uint64]common.Hash
	head          int64
	height        uint64
}

// NewEthereumMock returns a new EthereumMock client. It is used for testing.
func NewEthereumMock(genesis *canonicalStateChainContract.CanonicalStateChainHeader) *ethereumMock {

	e := &ethereumMock{
		rollupHeaders: make(map[common.Hash]canonicalStateChainContract.CanonicalStateChainHeader),
		indexToHash:   make(map[uint64]common.Hash),
		head:          -1,
	}

	e.PushRollupHead(genesis)
	return e
}

func (e *ethereumMock) Wait(txHash common.Hash) (*types.Receipt, error) {
	return types.NewReceipt(txHash[:], false, 21000), nil
}

// GetRollupHead returns the latest rollup block header.
func (e *ethereumMock) GetRollupHead() (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return e.GetRollupHeader(uint64(e.head))
}

// PushRollupHead pushes a new rollup block header.
func (e *ethereumMock) PushRollupHead(header *canonicalStateChainContract.CanonicalStateChainHeader) (*types.Transaction, error) {
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
func (e *ethereumMock) GetRollupHeader(index uint64) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	hash, ok := e.indexToHash[index]
	if !ok {
		return canonicalStateChainContract.CanonicalStateChainHeader{}, nil
	}

	return e.rollupHeaders[hash], nil
}

// GetRollupHeaderByHash returns the rollup block header with the given hash.
func (e *ethereumMock) GetRollupHeaderByHash(hash common.Hash) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
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
