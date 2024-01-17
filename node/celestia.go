package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"hummingbird/node/contracts"
	"log/slog"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/ethereum/go-ethereum/common"

	"github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/celestiaorg/celestia-app/pkg/square"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	blobstreamtypes "github.com/celestiaorg/celestia-app/x/qgb/types"
)

// CelestiaPointer is a pointer to a Celestia header
type CelestiaPointer struct {
	Height   uint64
	DataRoot common.Hash
	TxHash   common.Hash
}

type CelestiaProof struct {
	Nonce        *big.Int
	Tuple        *contracts.DataRootTuple
	WrappedProof *contracts.BinaryMerkleProof
}

// Celestia is the interface for interacting with the Celestia node
type Celestia interface {
	Namespace() string
	PublishBundle(blocks Bundle) (*CelestiaPointer, error)
	GetProof(txHash []byte) (*CelestiaProof, error)
}

type CelestiaClientOpts struct {
	Endpoint      string
	Token         string
	TendermintRPC string
	GRPC          string
	Namespace     string
	Logger        *slog.Logger
	GasPrice      float64
}

type CelestiaClient struct {
	namespace string
	client    *client.Client
	trpc      *http.HTTP
	grcp      *grpc.ClientConn
	logger    *slog.Logger
	gasPrice  float64
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

	// post the blob
	pointer, err := c.submitBlob(context.Background(), cosmosmath.NewInt(fee), gasLimit, []*blob.Blob{b})
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

	pointer := &CelestiaPointer{
		Height:   uint64(response.Height),
		DataRoot: common.BytesToHash(blobs[0].Commitment),
		TxHash:   common.BytesToHash(txHash),
	}

	return pointer, err
}

func (c *CelestiaClient) GetProof(txHash []byte) (*CelestiaProof, error) {
	ctx := context.Background()

	// Get the tx
	tx, err := c.trpc.Tx(ctx, txHash, true)
	if err != nil {
		return nil, err
	}

	// Get the block that contains the tx
	blockRes, err := c.trpc.Block(context.Background(), &tx.Height)
	if err != nil {
		return nil, err
	}

	// Get the tx share range inside the block
	shareRange, err := square.TxShareRange(blockRes.Block.Data.Txs.ToSliceOfBytes(), int(tx.Index), blockRes.Block.Header.Version.App)
	if err != nil {
		return nil, err
	}

	// Get the shares proof
	sharesProofs, err := c.trpc.ProveShares(ctx, uint64(tx.Height), uint64(shareRange.Start), uint64(shareRange.End))
	if err != nil {
		return nil, err
	}

	// Verify the shares proof
	if !sharesProofs.VerifyProof() {
		return nil, err
	}

	// New gRPC query client
	queryClient := blobstreamtypes.NewQueryClient(c.grcp)

	// Get the data commitment range for the block height
	resp, err := queryClient.DataCommitmentRangeForHeight(ctx, &blobstreamtypes.QueryDataCommitmentRangeForHeightRequest{Height: uint64(tx.Height)})
	if err != nil {
		return nil, err
	}

	// Get the data root inclusion proof
	dcProof, err := c.trpc.DataRootInclusionProof(ctx, uint64(tx.Height), resp.DataCommitment.BeginBlock, resp.DataCommitment.EndBlock)
	if err != nil {
		return nil, err
	}

	tuple := contracts.DataRootTuple{
		Height:   big.NewInt(int64(tx.Height)),
		DataRoot: *(*[32]byte)(blockRes.Block.DataHash),
	}

	sideNodes := make([][32]byte, len(dcProof.Proof.Aunts))
	for i, aunt := range dcProof.Proof.Aunts {
		sideNodes[i] = *(*[32]byte)(aunt)
	}
	wrappedProof := contracts.BinaryMerkleProof{
		SideNodes: sideNodes,
		Key:       big.NewInt(dcProof.Proof.Index),
		NumLeaves: big.NewInt(dcProof.Proof.Total),
	}

	proof := &CelestiaProof{
		Nonce:        big.NewInt(int64(resp.DataCommitment.Nonce)),
		Tuple:        &tuple,
		WrappedProof: &wrappedProof,
	}

	return proof, nil
}

// MOCK CLINT FOR TESTING

type celestiaMock struct {
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

func (c *celestiaMock) Namespace() string {
	return c.namespace
}

func (c *celestiaMock) PublishBundle(blocks Bundle) (*CelestiaPointer, error) {
	c.height++

	// use the first block's hash as the data root
	hash := blocks.Blocks[0].Hash()
	c.blocks[hash] = blocks

	c.pointers[hash] = &CelestiaPointer{
		Height:   c.height,
		DataRoot: hash,
		TxHash:   hash,
	}

	return c.pointers[hash], nil
}

// returns a mock proof, cannot be used for verification
func (c *celestiaMock) GetProof(hash []byte) (*CelestiaProof, error) {
	_, ok := c.blocks[common.BytesToHash(hash)]
	if !ok {
		return nil, blob.ErrBlobNotFound
	}

	p, ok := c.pointers[common.BytesToHash(hash)]
	if !ok {
		return nil, blob.ErrBlobNotFound
	}

	return &CelestiaProof{
		Nonce: big.NewInt(0),
		Tuple: &contracts.DataRootTuple{
			Height:   big.NewInt(int64(p.Height)),
			DataRoot: p.DataRoot,
		},
		WrappedProof: &contracts.BinaryMerkleProof{
			SideNodes: make([][32]byte, 0),
			Key:       big.NewInt(0),
			NumLeaves: big.NewInt(0),
		},
	}, nil
}
