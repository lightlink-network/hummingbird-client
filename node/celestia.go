package node

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/state"

	// "github.com/celestiaorg/celestia-node/share"
	// "github.com/celestiaorg/celestia-openrpc/types/share"
	openclient "github.com/celestiaorg/celestia-openrpc"
	gosquare "github.com/celestiaorg/go-square/square"
	"github.com/celestiaorg/go-square/v2/share"
	"github.com/ethereum/go-ethereum/common"
	thttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"

	"github.com/celestiaorg/celestia-app/v4/pkg/appconsts"
	blobtypes "github.com/celestiaorg/celestia-app/v4/x/blob/types"
	"github.com/celestiaorg/go-square/shares"

	blobstreamXContract "hummingbird/node/contracts/BlobstreamX.sol"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
	"hummingbird/utils"
)

// CelestiaPointer is a pointer to a Celestia header
type CelestiaPointer struct {
	Height     uint64
	ShareStart uint64
	ShareLen   uint64

	// Extra data Only present if the pointer is stored in the local database.
	Commitment common.Hash
	TxHash     common.Hash
}

type CelestiaProof struct {
	Nonce        *big.Int
	Tuple        *blobstreamXContract.DataRootTuple
	WrappedProof *challengeContract.BinaryMerkleProof
}

// Celestia is the interface for interacting with the Celestia node
type Celestia interface {
	Namespace() string
	PublishBundle(blocks Bundle) (*CelestiaPointer, float64, error)
	GetProof(pointer *CelestiaPointer, startBlock uint64, endBlock uint64, proofNonce big.Int) (*CelestiaProof, error)
	GetSharesByNamespace(pointer *CelestiaPointer) ([]shares.Share, error)
	GetSharesByPointer(pointer *CelestiaPointer) ([]shares.Share, error)
	GetShareProof(celestiaPointer *CelestiaPointer, shareIndex uint32) (*types.ShareProof, error)
	GetSharesProof(celestiaPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error)
	GetPointer(txHash common.Hash) (*CelestiaPointer, error)
}

type CelestiaClientOpts struct {
	Endpoint                string
	Token                   string
	TendermintRPC           string
	Namespace               string
	Logger                  *slog.Logger
	GasPrice                float64
	GasPriceIncreasePercent *big.Int
	GasAPI                  string
	Retries                 int
	RetryDelay              time.Duration
}

var _ Celestia = &CelestiaClient{}

type CelestiaClient struct {
	namespace               string
	client                  *client.Client
	openrpcClient           *openclient.Client
	trpc                    *thttp.HTTP
	logger                  *slog.Logger
	gasPrice                float64
	gasPriceIncreasePercent *big.Int
	gasAPI                  string
	retries                 int
	retryDelay              time.Duration
}

func NewCelestiaClient(opts CelestiaClientOpts) (*CelestiaClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	c, err := client.NewClient(context.Background(), opts.Endpoint, opts.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Celestia: %w", err)
	}

	openrpcClient, err := openclient.NewClient(context.Background(), opts.Endpoint, opts.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Celestia OpenRPC: %w", err)
	}

	trpc, err := thttp.New(opts.TendermintRPC, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tendermint RPC: %w", err)
	}

	if err := trpc.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Tendermint RPC: %w", err)
	}

	opts.Logger.Info("Connected to Celestia")
	return &CelestiaClient{
		namespace:               opts.Namespace,
		client:                  c,
		openrpcClient:           openrpcClient,
		trpc:                    trpc,
		logger:                  opts.Logger,
		gasPrice:                opts.GasPrice,
		gasPriceIncreasePercent: opts.GasPriceIncreasePercent,
		gasAPI:                  opts.GasAPI,
		retries:                 opts.Retries,
		retryDelay:              opts.RetryDelay,
	}, nil
}

func (c *CelestiaClient) Namespace() string {
	return c.namespace
}

func (c *CelestiaClient) PublishBundle(blocks Bundle) (*CelestiaPointer, float64, error) {
	// get the namespace
	ns, err := share.NewV0Namespace([]byte(c.Namespace()))
	if err != nil {
		return nil, 0, err
	}

	// encode the blocks
	enc, err := blocks.EncodeRLP()
	if err != nil {
		return nil, 0, err
	}

	// create blob to submit
	//b, err := blob.NewBlob(0, ns, []byte(enc))
	b, err := blob.NewBlobV0(ns, []byte(enc))
	if err != nil {
		panic(err)
	}

	// gas price is defined by each node operator. 0.003 is a good default to be accepted
	gasPrice := c.GasPrice()

	if c.gasPriceIncreasePercent != nil {
		apiPrice := gasPrice
		gasPrice *= 1 + float64(c.gasPriceIncreasePercent.Int64())/100
		c.logger.Info("Gas price increased", "percent", c.gasPriceIncreasePercent, "old_gas_price", apiPrice, "new_gas_price", gasPrice)
	}

	// estimate gas limit (maximum gas used by the tx)
	gasLimit := blobtypes.DefaultEstimateGas([]uint32{uint32(b.DataLen())})

	var pointer *CelestiaPointer

	i := 0
	for {
		// post the blob
		pointer, err = c.submitBlob(context.Background(), gasPrice, gasLimit, []*blob.Blob{b})
		if err == nil || i >= c.retries {
			break
		}

		// Increase gas price by 20% if the transaction fails
		gasPrice *= 1.2

		c.logger.Warn("Failed to submit blob, retrying after delay", "delay", c.retryDelay, "attempt", i+1, "gas_limit", gasLimit, "gas_price", gasPrice, "error", err)

		i++

		// Delay between publishing bundles to Celestia to mitigate 'incorrect account sequence' errors
		time.Sleep(c.retryDelay)
	}

	if err != nil {
		return nil, gasPrice, err
	}

	return pointer, gasPrice, nil
}

// PostData submits a new transaction with the provided data to the Celestia node.
func (c *CelestiaClient) submitBlob(ctx context.Context, gasPrice float64, gasLimit uint64, blobs []*blob.Blob) (*CelestiaPointer, error) {
	//response, err := c.client.State.SubmitPayForBlob(ctx, fee, gasLimit, blobs)
	response, err := c.client.State.SubmitPayForBlob(ctx, blob.ToLibBlobs(blobs...), state.NewTxConfig(
		state.WithGas(gasLimit),
		state.WithGasPrice(gasPrice),
	))
	if err != nil {
		return nil, err
	}

	txHash, err := hex.DecodeString(response.TxHash)
	if err != nil {
		return nil, err
	}

	// Delay here before getting the block to ensure the tx is included
	time.Sleep(5 * time.Second)

	// Get the block that contains the tx
	pointer, err := c.GetPointer(common.BytesToHash(txHash))
	if err != nil {
		return nil, err
	}

	return pointer, err
}

func (c *CelestiaClient) GetProof(pointer *CelestiaPointer, startBlock uint64, endBlock uint64, proofNonce big.Int) (*CelestiaProof, error) {
	ctx := context.Background()

	blockHeight := int64(pointer.Height)

	// Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &blockHeight)
	if err != nil {
		return nil, err
	}

	// Get the shares proof
	sharesProofs, err := c.trpc.ProveShares(ctx, pointer.Height, pointer.ShareStart, pointer.ShareStart+pointer.ShareLen)
	if err != nil {
		return nil, err
	}

	// Verify the shares proof
	if !sharesProofs.VerifyProof() {
		return nil, err
	}

	// Get the data root inclusion proof
	dcProof, err := c.trpc.DataRootInclusionProof(ctx, uint64(blockHeight), startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	tuple := blobstreamXContract.DataRootTuple{
		Height:   big.NewInt(blockHeight),
		DataRoot: *(*[32]byte)(blockRes.Block.DataHash),
	}

	sideNodes := make([][32]byte, len(dcProof.Proof.Aunts))
	for i, aunt := range dcProof.Proof.Aunts {
		sideNodes[i] = *(*[32]byte)(aunt)
	}
	wrappedProof := challengeContract.BinaryMerkleProof{
		SideNodes: sideNodes,
		Key:       big.NewInt(dcProof.Proof.Index),
		NumLeaves: big.NewInt(dcProof.Proof.Total),
	}

	proof := &CelestiaProof{
		Nonce:        &proofNonce,
		Tuple:        &tuple,
		WrappedProof: &wrappedProof,
	}

	return proof, nil
}

// GetPointer returns the pointer to the Celestia header that contains the tx with the given hash
func (c *CelestiaClient) GetPointer(txHash common.Hash) (*CelestiaPointer, error) {
	tx, err := c.trpc.Tx(context.Background(), txHash.Bytes(), true)
	if err != nil {
		return nil, err
	}
	// Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &tx.Height)
	if err != nil {
		return nil, err
	}
	// Get the blob share range inside the block, using square instead

	maxSquareSize := appconsts.SquareSizeUpperBound
	subtreeRootThreshold := appconsts.SubtreeRootThreshold
	blobShareRange, err := gosquare.BlobShareRange(blockRes.Block.Txs.ToSliceOfBytes(), int(tx.Index), int(0), maxSquareSize, subtreeRootThreshold)
	if err != nil {
		return nil, err
	}
	return &CelestiaPointer{
		Height:     uint64(tx.Height),
		Commitment: common.BytesToHash(blockRes.Block.DataHash),
		ShareStart: uint64(blobShareRange.Start),
		ShareLen:   uint64(blobShareRange.End - blobShareRange.Start),
		TxHash:     txHash,
	}, nil
}

func (c *CelestiaClient) GetSharesByNamespace(pointer *CelestiaPointer) ([]shares.Share, error) {
	ctx := context.Background()

	// 1. Namespace
	// ns, err := openshare.NewBlobNamespaceV0([]byte(c.Namespace()))
	// if err != nil {
	// 	return nil, fmt.Errorf("GetShares: failed to get namespace: %w", err)
	// }
	ns, err := share.NewV0Namespace([]byte(c.Namespace()))
	if err != nil {
		return nil, fmt.Errorf("GetShares: failed to get namespace: %w", err)
	}

	// 3. Get the shares
	//s, err := c.client.Share.GetSharesByNamespace(ctx, h, ns)
	nsData, err := c.client.Share.GetNamespaceData(ctx, pointer.Height, ns)
	if err != nil {
		return nil, fmt.Errorf("GetShares: failed to get namespace data: %w", err)
	}

	// s, err := c.openrpcClient.Share.GetSharesByNamespace(ctx, h, ns)
	// if err != nil {
	// 	return nil, fmt.Errorf("GetShares: failed to get shares: %w", err)
	// }

	return utils.NSSharesToShares(nsData.Flatten()), nil
}

func (c *CelestiaClient) GetSharesByPointer(pointer *CelestiaPointer) ([]shares.Share, error) {
	ctx := context.Background()

	proof, err := c.trpc.ProveShares(ctx, pointer.Height, pointer.ShareStart, pointer.ShareStart+pointer.ShareLen)
	if err != nil {
		return nil, err
	}

	return utils.BytesToShares(proof.Data)
}

func (c *CelestiaClient) GetShareProof(celestiaPointer *CelestiaPointer, shareIndex uint32) (*types.ShareProof, error) {
	ctx := context.Background()

	shareStart := celestiaPointer.ShareStart + uint64(shareIndex)
	shareEnd := celestiaPointer.ShareStart + uint64(shareIndex+1)

	// Get the shares proof
	sharesProofs, err := c.trpc.ProveShares(ctx, celestiaPointer.Height, shareStart, shareEnd)
	if err != nil {
		return nil, err
	}

	// Verify the shares proof
	if !sharesProofs.VerifyProof() {
		return nil, err
	}

	return &sharesProofs, nil
}

func (c *CelestiaClient) GetSharesProof(celPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error) {
	ctx := context.Background()

	shareStart := celPointer.ShareStart + uint64(sharePointer.StartShare)
	shareEnd := celPointer.ShareStart + uint64(sharePointer.EndShare()+1)

	// Get the shares proof
	sharesProofs, err := c.trpc.ProveShares(ctx, celPointer.Height, shareStart, shareEnd)
	if err != nil {
		return nil, err
	}

	// Verify the shares proof
	if !sharesProofs.VerifyProof() {
		return nil, err
	}

	return &sharesProofs, nil
}

type GasPrice struct {
	Slow   string `json:"slow"`
	Median string `json:"median"`
	Fast   string `json:"fast"`
}

func (c *CelestiaClient) GasPrice() float64 {
	// Make HTTP GET request
	resp, err := http.Get(c.gasAPI)
	if err != nil {
		c.logger.Error("Error making HTTP request", "error", err)
		return c.gasPrice
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Error reading response body", "error", err)
		return c.gasPrice
	}

	// Parse JSON response
	var gasPrice GasPrice
	err = json.Unmarshal(body, &gasPrice)
	if err != nil {
		c.logger.Error("Error parsing JSON response", "error", err)
		return c.gasPrice
	}

	// Convert fast gas price to float64
	fast, err := strconv.ParseFloat(gasPrice.Fast, 64)
	if err != nil {
		c.logger.Error("Error converting fast gas price to float64", "error", err)
		return c.gasPrice
	}

	return fast
}

// MOCK CLINT FOR TESTING

var _ Celestia = &celestiaMock{}

type celestiaMock struct {
	fakeProof bool
	namespace string
	height    uint64
	blocks    map[common.Hash]Bundle
	pointers  map[common.Hash]*CelestiaPointer
}

// NewCelestiaMock returns a new CelestiaMock client. It is used for testing.
func NewCelestiaMock(namespace string) *celestiaMock {
	return &celestiaMock{
		namespace: namespace,
		blocks:    make(map[common.Hash]Bundle),
		pointers:  make(map[common.Hash]*CelestiaPointer),
	}
}

func (c *celestiaMock) SetFakeProof(b bool) {
	c.fakeProof = b
}

func (c *celestiaMock) Namespace() string {
	return c.namespace
}

func (c *celestiaMock) PublishBundle(blocks Bundle) (*CelestiaPointer, float64, error) {
	c.height++

	// use the first block's hash as the data root
	hash := blocks.Blocks[0].Hash()
	c.blocks[hash] = blocks

	c.pointers[hash] = &CelestiaPointer{
		Height:     c.height,
		ShareStart: 0,
		ShareLen:   uint64(len(blocks.Blocks)),
		TxHash:     hash,
	}

	return c.pointers[hash], 0, nil
}

// returns a mock proof, cannot be used for verification
func (c *celestiaMock) GetProof(pointer *CelestiaPointer, startBlock uint64, endBlock uint64, proofNonce big.Int) (*CelestiaProof, error) {
	if !c.fakeProof {
		return nil, fmt.Errorf("failed")
	}
	return &CelestiaProof{
		Nonce: big.NewInt(0),
		Tuple: &blobstreamXContract.DataRootTuple{
			Height:   new(big.Int).SetUint64(pointer.Height),
			DataRoot: pointer.Commitment,
		},
		WrappedProof: &challengeContract.BinaryMerkleProof{
			SideNodes: make([][32]byte, 0),
			Key:       big.NewInt(0),
			NumLeaves: big.NewInt(0),
		},
	}, nil
}

func (c *celestiaMock) GetSharesByNamespace(pointer *CelestiaPointer) ([]shares.Share, error) {
	return nil, nil
}

func (c *celestiaMock) GetSharesProof(celestiaPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error) {
	return nil, nil
}

func (c *celestiaMock) GetShareProof(celestiaPointer *CelestiaPointer, shareIndex uint32) (*types.ShareProof, error) {
	return nil, nil
}

func (c *celestiaMock) GetPointer(txHash common.Hash) (*CelestiaPointer, error) {
	return c.pointers[txHash], nil
}

func (c *celestiaMock) GetSharesByPointer(pointer *CelestiaPointer) ([]shares.Share, error) {
	return nil, nil
}
