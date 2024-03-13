package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	cosmosmath "cosmossdk.io/math"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	gosquare "github.com/celestiaorg/go-square/square"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	blobstreamtypes "github.com/celestiaorg/celestia-app/x/qgb/types"

	challengeContract "hummingbird/node/contracts/Challenge.sol"
	daOracleContract "hummingbird/node/contracts/DAOracle.sol"
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
	Tuple        *daOracleContract.DataRootTuple
	WrappedProof *challengeContract.BinaryMerkleProof
}

// Celestia is the interface for interacting with the Celestia node
type Celestia interface {
	Namespace() string
	PublishBundle(blocks Bundle) (*CelestiaPointer, error)
	GetProof(pointer *CelestiaPointer) (*CelestiaProof, error)
	GetShares(pointer *CelestiaPointer) ([]shares.Share, error)
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
	}

	if err != nil {
		return nil, err
	}

	return pointer, nil
}

func (c *CelestiaClient) waitForTxInclusion(ctx context.Context, h common.Hash) (*ctypes.ResultTx, error) {
	txHash := []byte{}
	maxRetries := 1000
	retryInterval := 10 * time.Second

	c.logger.Debug("Waiting for tx inclusion", "blob_hash", h.Hex())

	// Scan the mempool every 'retryInterval' until 'maxRetries' is reached or the tx is not found in the mempool
	// Scanning the pool gives us the tx hash
	for i := 0; i < maxRetries; i++ {
		txns, err := c.trpc.UnconfirmedTxs(ctx, nil)
		if err != nil {
			return nil, err
		}

		found := false
		for _, tx := range txns.Txs {
			blobtx, isBlob := types.UnmarshalBlobTx(tx)
			if !isBlob {
				c.logger.Info("waitForTxInclusion: discovered tx is not a blob")
				continue
			}
			if len(blobtx.Blobs) == 0 {
				return nil, fmt.Errorf("waitForTxInclusion: discovered tx has no blobs")
			}
			// Check if the data hash of the tx matches the hash of the data we submitted
			if crypto.Keccak256Hash(blobtx.Blobs[0].Data) == h {
				txHash = tx.Hash()
				c.logger.Debug("Tx found in mempool", "tx", h.Hex())
				found = true
				break // if the tx is found in the mempool, break the inner loop as it is not included in a block yet
			}
		}

		if !found {
			break // if the tx is not found in the mempool, break the outer loop as it is already included in a block
		}

		time.Sleep(retryInterval)
	}

	if len(txHash) == 0 {
		return nil, fmt.Errorf("waitForTxInclusion: tx not found after max retries:  %v", maxRetries)
	}

	// Get the tx with block inclusion info
	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval)

		tx, err := c.trpc.Tx(ctx, txHash, true)
		if err != nil {
			// Sometimes the tx is not in the mempool, but not found in the block yet
			// We need to wait for it to be included in a block
			if strings.Contains(err.Error(), "not found") {
				continue
			} else {
				return nil, err
			}
		}
		if tx.Height == 0 {
			c.logger.Debug("Tx found but height is 0, retrying...", "tx", tx.Hash, "height", tx.Height, "index", tx.Index)
			time.Sleep(retryInterval)
			continue
		}
		c.logger.Debug("Tx found and included in a block", "tx", tx.Hash, "height", tx.Height, "index", tx.Index)
		return tx, nil
	}

	return nil, fmt.Errorf("waitForTxInclusion: tx not found after max retries:  %v", maxRetries)
}

// PostData submits a new transaction with the provided data to the Celestia node.
func (c *CelestiaClient) submitBlob(ctx context.Context, fee cosmosmath.Int, gasLimit uint64, blobs []*blob.Blob) (*CelestiaPointer, error) {
	tx := &ctypes.ResultTx{}
	response, err := c.client.State.SubmitPayForBlob(ctx, fee, gasLimit, blobs)
	if err != nil {
		if strings.Contains(err.Error(), "timed out waiting for tx to be included in a block") {
			tx, err = c.waitForTxInclusion(ctx, crypto.Keccak256Hash(blobs[0].Data))
			if err != nil {
				return nil, fmt.Errorf("submitBlob: failed to wait for tx inclusion: %w", err)
			}
			if tx == nil {
				return nil, fmt.Errorf("submitBlob: failed to wait for tx inclusion: %w", err)
			}
		} else {
			return nil, err
		}
	}

	var txHash []byte

	if response != nil {
		txHash, err = hex.DecodeString(response.TxHash)
		if err != nil {
			return nil, err
		}
	} else {
		txHash = tx.Hash
	}

	// Get the block that contains the tx
	pointer, err := c.GetPointer(common.BytesToHash(txHash))
	if err != nil {
		return nil, err
	}

	return pointer, err
}

func (c *CelestiaClient) GetProof(pointer *CelestiaPointer) (*CelestiaProof, error) {
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

	// New gRPC query client
	queryClient := blobstreamtypes.NewQueryClient(c.grcp)

	// Get the data commitment range for the block height
	resp, err := queryClient.DataCommitmentRangeForHeight(ctx, &blobstreamtypes.QueryDataCommitmentRangeForHeightRequest{Height: uint64(blockHeight)})
	if err != nil {
		return nil, err
	}

	// Get the data root inclusion proof
	dcProof, err := c.trpc.DataRootInclusionProof(ctx, uint64(blockHeight), resp.DataCommitment.BeginBlock, resp.DataCommitment.EndBlock)
	if err != nil {
		return nil, err
	}

	tuple := daOracleContract.DataRootTuple{
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
		Nonce:        big.NewInt(int64(resp.DataCommitment.Nonce)),
		Tuple:        &tuple,
		WrappedProof: &wrappedProof,
	}

	return proof, nil
}

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

func (c *CelestiaClient) GetShares(pointer *CelestiaPointer) ([]shares.Share, error) {
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
func (c *celestiaMock) GetProof(pointer *CelestiaPointer) (*CelestiaProof, error) {
	if !c.fakeProof {
		return nil, fmt.Errorf("failed")
	}
	return &CelestiaProof{
		Nonce: big.NewInt(0),
		Tuple: &daOracleContract.DataRootTuple{
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

func (c *celestiaMock) GetShares(pointer *CelestiaPointer) ([]shares.Share, error) {
	return nil, nil
}

func (c *celestiaMock) GetSharesProof(celestiaPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error) {
	return nil, nil
}

func (c *celestiaMock) GetPointer(txHash common.Hash) (*CelestiaPointer, error) {
	return c.pointers[txHash], nil
}
