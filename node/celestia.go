package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	cosmosmath "cosmossdk.io/math"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	gosquare "github.com/celestiaorg/go-square/square"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"

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
	PublishBundle(blocks Bundle) (*CelestiaPointer, error)
	GetProof(pointer *CelestiaPointer, startBlock uint64, endBlock uint64, proofNonce big.Int) (*CelestiaProof, error)
	GetSharesByNamespace(pointer *CelestiaPointer) ([]shares.Share, error)
	GetSharesByPointer(pointer *CelestiaPointer) ([]shares.Share, error)
	GetShareProof(celestiaPointer *CelestiaPointer, shareIndex uint32) (*types.ShareProof, error)
	GetSharesProof(celestiaPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error)
	GetPointer(txHash common.Hash) (*CelestiaPointer, error)
}

type CelestiaClientOpts struct {
	Endpoint      string
	Token         string
	TendermintRPC string
	GRPC          string
	Namespace     string
	Logger        *slog.Logger
	GasPrice      float64
	Retries       int
}

type CelestiaClient struct {
	namespace string
	client    *client.Client
	trpc      *http.HTTP
	grcp      *grpc.ClientConn
	logger    *slog.Logger
	gasPrice  float64
	retries   int
}

func NewCelestiaClient(opts CelestiaClientOpts) (*CelestiaClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	c, err := client.NewClient(context.Background(), opts.Endpoint, opts.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Celestia: %w", err)
	}

	trpc, err := http.New(opts.TendermintRPC, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tendermint RPC: %w", err)
	}

	if err := trpc.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Tendermint RPC: %w", err)
	}

	grcp, err := grpc.Dial(opts.GRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Celestia GRPC: %w", err)
	}

	opts.Logger.Info("Connected to Celestia")
	return &CelestiaClient{
		namespace: opts.Namespace,
		client:    c,
		trpc:      trpc,
		grcp:      grcp,
		logger:    opts.Logger,
		gasPrice:  opts.GasPrice,
		retries:   opts.Retries,
	}, nil
}

func (c *CelestiaClient) Namespace() string {
	return c.namespace
}

func (c *CelestiaClient) PublishBundle(blocks Bundle) (*CelestiaPointer, error) {
	// get the namespace
	ns, err := share.NewBlobNamespaceV0([]byte(c.Namespace()))
	if err != nil {
		return nil, err
	}

	// encode the blocks
	enc, err := blocks.EncodeRLP()
	if err != nil {
		return nil, err
	}

	// create blob to submit
	b, err := blob.NewBlob(0, ns, []byte(enc))
	if err != nil {
		panic(err)
	}

	// gas price is defined by each node operator. 0.003 is a good default to be accepted
	gasPrice := c.gasPrice

	// estimate gas limit (maximum gas used by the tx)
	gasLimit := blobtypes.DefaultEstimateGas([]uint32{uint32(b.Size())})

	// fee is gas price * gas limit. State machine does not refund users for unused gas so all of the fee is used
	fee := int64(gasPrice * float64(gasLimit))

	var pointer *CelestiaPointer

	i := 0
	for {
		// post the blob
		pointer, err = c.submitBlob(context.Background(), cosmosmath.NewInt(fee), gasLimit, []*blob.Blob{b})
		if err == nil || i >= c.retries {
			break
		}

		// Increase gas price by 20% if the transaction fails
		gasPrice *= 1.2
		fee = int64(gasPrice * float64(gasLimit))

		c.logger.Warn("Failed to submit blob, retrying", "attempt", i+1, "fee", fee, "gas_limit", gasLimit, "gas_price", gasPrice, "error", err)

		i++

		// Delay between publishing bundles to Celestia to mitigate 'incorrect account sequence' errors
		time.Sleep(30 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return pointer, nil
}

// PostData submits a new transaction with the provided data to the Celestia node.
func (c *CelestiaClient) submitBlob(ctx context.Context, fee cosmosmath.Int, gasLimit uint64, blobs []*blob.Blob) (*CelestiaPointer, error) {
	response, err := c.client.State.SubmitPayForBlob(ctx, fee, gasLimit, blobs)
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
	version := blockRes.Block.Header.Version.App
	maxSquareSize := appconsts.SquareSizeUpperBound(version)
	subtreeRootThreshold := appconsts.SubtreeRootThreshold(version)
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
	ns, err := share.NewBlobNamespaceV0([]byte(c.Namespace()))
	if err != nil {
		return nil, fmt.Errorf("GetShares: failed to get namespace: %w", err)
	}

	// 0. Get the header
	h, err := c.client.Header.GetByHeight(ctx, pointer.Height)
	if err != nil {
		return nil, fmt.Errorf("GetShares: failed to get header: %d %w", pointer.Height, err)
	}

	// 3. Get the shares
	s, err := c.client.Share.GetSharesByNamespace(ctx, h, ns)
	if err != nil {
		return nil, fmt.Errorf("GetShares: failed to get shares: %w", err)
	}

	return utils.NSSharesToShares(s), nil
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

// MOCK CLINT FOR TESTING

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

func (c *celestiaMock) PublishBundle(blocks Bundle) (*CelestiaPointer, error) {
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

	return c.pointers[hash], nil
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
