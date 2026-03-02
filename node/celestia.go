package node

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/state"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	gosquare "github.com/celestiaorg/go-square/v3"
	"github.com/celestiaorg/go-square/v3/share"
	thttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/celestiaorg/celestia-app/v6/app"
	"github.com/celestiaorg/celestia-app/v6/app/encoding"
	"github.com/celestiaorg/celestia-app/v6/pkg/appconsts"

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
	GetSharesByNamespace(pointer *CelestiaPointer) ([]share.Share, error)
	GetSharesByPointer(pointer *CelestiaPointer) ([]share.Share, error)
	GetShareProof(celestiaPointer *CelestiaPointer, shareIndex uint32) (*types.ShareProof, error)
	GetSharesProof(celestiaPointer *CelestiaPointer, sharePointer *SharePointer) (*types.ShareProof, error)
	GetPointer(txHash common.Hash) (*CelestiaPointer, error)
}

type CelestiaClientOpts struct {
	ConsensusRPC  string
	Namespace     string
	Logger        *slog.Logger
	GasPrice      float64
	Retries       int
	RetryDelay    time.Duration
	Mnemonic      string
	ConsensusGRPC string
	ConsensusTLS  bool
	Network       string
}

const celestiaKeyName = "hummingbird"

var _ Celestia = &CelestiaClient{}

type CelestiaClient struct {
	namespace  string
	coreAccess *state.CoreAccessor
	trpc       *thttp.HTTP
	logger     *slog.Logger
	gasPrice   float64
	retries    int
	retryDelay time.Duration
}

func NewCelestiaClient(opts CelestiaClientOpts) (*CelestiaClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	trpc, err := thttp.New(opts.ConsensusRPC, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to consensus RPC: %w", err)
	}

	if err := trpc.Start(); err != nil {
		return nil, fmt.Errorf("failed to start consensus RPC: %w", err)
	}

	encCfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	kr := keyring.NewInMemory(encCfg.Codec)

	_, err = kr.NewAccount(
		celestiaKeyName,
		opts.Mnemonic,
		keyring.DefaultBIP39Passphrase,
		sdk.GetConfig().GetFullBIP44Path(),
		hd.Secp256k1,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to import celestia key from mnemonic: %w", err)
	}

	rec, err := kr.Key(celestiaKeyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get celestia key: %w", err)
	}
	addr, err := rec.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get celestia address: %w", err)
	}
	opts.Logger.Info("Celestia signer address", "address", addr.String())

	var grpcOpts []grpc.DialOption
	if opts.ConsensusTLS {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	grpcConn, err := grpc.NewClient(opts.ConsensusGRPC, grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to celestia gRPC at %s: %w", opts.ConsensusGRPC, err)
	}

	ctx := context.Background()
	coreAccessor, err := state.NewCoreAccessor(kr, celestiaKeyName, nil, grpcConn, opts.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to create celestia core accessor: %w", err)
	}

	if err := coreAccessor.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start celestia core accessor: %w", err)
	}

	opts.Logger.Info("Connected to Celestia", "grpc", opts.ConsensusGRPC, "network", opts.Network)
	return &CelestiaClient{
		namespace:  opts.Namespace,
		coreAccess: coreAccessor,
		trpc:       trpc,
		logger:     opts.Logger,
		gasPrice:   opts.GasPrice,
		retries:    opts.Retries,
		retryDelay: opts.RetryDelay,
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
	b, err := blob.NewBlobV0(ns, []byte(enc))
	if err != nil {
		panic(err)
	}

	var pointer *CelestiaPointer

	i := 0
	for {
		// post the blob
		pointer, err = c.submitBlob(context.Background(), []*blob.Blob{b})
		if err == nil || i >= c.retries {
			break
		}

		c.logger.Warn("Failed to submit blob, retrying after delay", "delay", c.retryDelay, "attempt", i+1, "error", err)

		i++

		// Delay between publishing bundles to Celestia to mitigate 'incorrect account sequence' errors
		time.Sleep(c.retryDelay)
	}

	if err != nil {
		return nil, 0, err
	}

	return pointer, 0, nil
}

// PostData submits a new transaction with the provided data to the Celestia node.
func (c *CelestiaClient) submitBlob(ctx context.Context, blobs []*blob.Blob) (*CelestiaPointer, error) {
	c.logger.Debug("Submitting blob to Celestia",
		"blob_count", len(blobs),
		"blob_sizes", func() []int {
			sizes := make([]int, len(blobs))
			for i, b := range blobs {
				sizes[i] = b.DataLen()
			}
			return sizes
		}())

	txConfig := state.NewTxConfig()

	c.logger.Debug("Calling SubmitPayForBlob",
		"endpoint", "State.SubmitPayForBlob",
		"tx_config", fmt.Sprintf("%+v", txConfig))

	response, err := c.coreAccess.SubmitPayForBlob(ctx, blob.ToLibBlobs(blobs...), txConfig)
	if err != nil {
		c.logger.Error("SubmitPayForBlob failed",
			"error", err,
			"error_type", fmt.Sprintf("%T", err))
		return nil, err
	}

	c.logger.Debug("SubmitPayForBlob response received",
		"tx_hash", response.TxHash,
		"height", response.Height,
		"response_type", fmt.Sprintf("%T", response))

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
	c.logger.Debug("GetPointer: Fetching transaction details",
		"tx_hash", txHash.Hex(),
		"tendermint_rpc_endpoint", c.trpc.Remote())

	tx, err := c.trpc.Tx(context.Background(), txHash.Bytes(), true)
	if err != nil {
		c.logger.Error("GetPointer: Failed to get transaction from Tendermint RPC",
			"tx_hash", txHash.Hex(),
			"error", err,
			"error_type", fmt.Sprintf("%T", err))
		return nil, err
	}

	c.logger.Debug("GetPointer: Transaction fetched successfully",
		"tx_hash", txHash.Hex(),
		"height", tx.Height,
		"index", tx.Index)
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
		c.logger.Error("GetPointer: Failed to get blob share range",
			"error", err,
			"error_type", fmt.Sprintf("%T", err))
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

func (c *CelestiaClient) GetSharesByNamespace(pointer *CelestiaPointer) ([]share.Share, error) {
	return c.GetSharesByPointer(pointer)
}

func (c *CelestiaClient) GetSharesByPointer(pointer *CelestiaPointer) ([]share.Share, error) {
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

func (c *celestiaMock) GetSharesByNamespace(pointer *CelestiaPointer) ([]share.Share, error) {
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

func (c *celestiaMock) GetSharesByPointer(pointer *CelestiaPointer) ([]share.Share, error) {
	return nil, nil
}
