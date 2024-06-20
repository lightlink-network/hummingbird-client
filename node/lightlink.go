package node

import (
	"encoding/json"
	"fmt"
	"hummingbird/node/jsonrpc"
	"hummingbird/node/lightlink"
	"hummingbird/utils"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// LinkLink is a client for the LightLink layer 2 network.
type LightLink interface {
	GetHeight() (uint64, error) // GetHeight returns the current height of the lightlink network.
	GetBlock(height uint64) (*types.Block, error)
	GetBlocks(start, end uint64) ([]*types.Block, error)
	GetOutputV0(last *types.Header) (OutputV0, error)
	GetProof(address common.Address, keys []string, height uint64) (*RawProof, error)
}

type LightLinkClientOpts struct {
	Endpoint                string
	Delay                   time.Duration
	Logger                  *slog.Logger
	L2ToL1MessagePasserAddr common.Address
}

type LightLinkClient struct {
	client *jsonrpc.Client
	opts   *LightLinkClientOpts
}

func NewLightLinkClient(opts *LightLinkClientOpts) (*LightLinkClient, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	client, err := jsonrpc.NewClient(opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LightLink: %w", err)
	}

	ll := &LightLinkClient{client: client, opts: opts}

	// check connection
	chainId, err := ll.GetChainId()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %w", err)
	}

	opts.Logger.Info("Connected to LightLink", "chainId", chainId)
	return ll, nil
}

func (l *LightLinkClient) GetChainId() (uint64, error) {
	resp, err := l.client.Call("eth_chainId", nil)
	if err != nil {
		return 0, err
	}

	numHex := resp.Result.(string)
	return hexutil.DecodeUint64(numHex)
}

func (l *LightLinkClient) GetHeight() (uint64, error) {
	resp, err := l.client.Call("eth_blockNumber", nil)
	if err != nil {
		return 0, err
	}

	numHex := resp.Result.(string)
	return hexutil.DecodeUint64(numHex)
}

func (l *LightLinkClient) GetBlock(height uint64) (*types.Block, error) {

	resp, err := l.client.Call("eth_getBlockByNumber", []any{hexutil.EncodeUint64(height), true})
	if err != nil {
		return nil, err
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response is not a map[string]interface{}: %v", resp.Result)
	}
	txs := types.Transactions{}
	for k, v := range result["transactions"].([]interface{}) {

		tx, err := unmarshalJsonTx(v)
		if err != nil {
			return nil, fmt.Errorf("failed to bind transaction %d: %w", k, err)
		}

		txs = append(txs, tx)
	}

	h := &types.Header{}
	err = resp.Bind(h)
	if err != nil {
		return nil, err
	}

	return types.NewBlockWithHeader(h).WithBody(txs, nil), nil
}

func (l *LightLinkClient) GetBlocks(start, end uint64) ([]*types.Block, error) {

	var blocks []*types.Block
	for i := start; i <= end; i++ {
		var block *types.Block
		var err error

		// retry up to 5 times in case of connreset or timeout errors etc
		for retry := 0; retry < 5; retry++ {
			block, err = l.GetBlock(i)
			if err == nil {
				break
			}
			time.Sleep(time.Second * time.Duration(2<<retry)) // exponential backoff
		}

		// if after 5 retries we still have an error, return it
		if err != nil {
			return nil, fmt.Errorf("failed to get block at height %d: %w", i, err)
		}

		// check if the block can be added to the bundle or if
		// the bundle has reached the max celestia tx size limit
		bundle := &Bundle{Blocks: append(blocks, block)}
		isUnderLimit, bundleSizeLimit, bundleEncodedSize, err := bundle.IsUnderTxLimit()
		if err != nil {
			return nil, fmt.Errorf("failed to check bundle size: %w", err)
		}

		if !isUnderLimit {
			l.opts.Logger.Info("Bundle has reached max celestia tx size limit", "blockCount", len(blocks), "bundleSize", bundleEncodedSize, "txSizeLimit", bundleSizeLimit)
			return blocks, nil
		}

		// add the block to the bundle
		blocks = append(blocks, block)

		// delay between requests
		if l.opts.Delay > 0 {
			time.Sleep(l.opts.Delay)
		}
	}

	return blocks, nil
}

func (l *LightLinkClient) GetWithdrawalRoot(height uint64) (common.Hash, error) {
	// get the storage root for L2ToL1MessagePasserAddr at the last block height
	proofRaw, err := l.client.Call("eth_getProof", []any{l.opts.L2ToL1MessagePasserAddr.Hex(), []string{}, hexutil.EncodeUint64(height)})
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get withdrawal address proof: %w", err)
	}
	proof := make(map[string]interface{})
	err = proofRaw.Bind(&proof)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to bind withdrawal address proof: %w", err)
	}
	withdrawalRoot := proof["storageHash"].(string)
	if withdrawalRoot == "" {
		return common.Hash{}, fmt.Errorf("failed to get withdrawal address proof: storageHash is empty")
	}

	return common.HexToHash(withdrawalRoot), nil
}

type RawProof struct {
	Address      string   `json:"address"`
	AccountProof []string `json:"accountProof"`
	Balance      string   `json:"balance"`
	CodeHash     string   `json:"codeHash"`
	Nonce        string   `json:"nonce"`
	StorageHash  string   `json:"storageHash"`
	StorageProof []struct {
		Key   string
		Value string
		Proof []string
	} `json:"storageProof"`
}

func (l *LightLinkClient) GetProof(address common.Address, keys []string, height uint64) (*RawProof, error) {
	proofRaw, err := l.client.Call("eth_getProof", []any{address.Hex(), keys, hexutil.EncodeUint64(height)})
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	proof := &RawProof{}
	err = proofRaw.Bind(proof)
	if err != nil {
		return nil, fmt.Errorf("failed to bind proof: %w", err)
	}

	return proof, nil
}

type OutputV0 struct {
	StateRoot                common.Hash
	MessagePasserStorageRoot common.Hash
	BlockHash                common.Hash
}

func (o OutputV0) Version() [32]byte {
	return [32]byte{}
}

func (o OutputV0) Root() common.Hash {
	var buf [128]byte
	version := o.Version()
	copy(buf[:32], version[:])
	copy(buf[32:], o.StateRoot[:])
	copy(buf[64:], o.MessagePasserStorageRoot[:])
	copy(buf[96:], o.BlockHash[:])
	return crypto.Keccak256Hash(buf[:])
}

func (l *LightLinkClient) GetOutputV0(last *types.Header) (OutputV0, error) {
	withdrawalRoot, err := l.GetWithdrawalRoot(last.Number.Uint64())
	if err != nil {
		return OutputV0{}, err
	}

	blockHash := utils.HashHeaderWithoutExtraData(last)
	return OutputV0{
		StateRoot:                last.Root,
		MessagePasserStorageRoot: withdrawalRoot,
		BlockHash:                blockHash,
	}, nil
}

func unmarshalJsonTx(from any) (*types.Transaction, error) {
	b, err := json.Marshal(from)
	if err != nil {
		return nil, err
	}

	return lightlink.UnMarshallTx(b)
}

// LightLinkMock is a mock LightLink client.

type lightLinkMock struct {
	Height uint64
	Blocks []*types.Block
}

func NewLightLinkMock() *lightLinkMock {
	return &lightLinkMock{Height: 0, Blocks: []*types.Block{}}
}

func (m *lightLinkMock) GetHeight() (uint64, error) {
	return m.Height, nil
}

func (m *lightLinkMock) GetBlock(height uint64) (*types.Block, error) {
	return m.Blocks[height], nil
}

func (m *lightLinkMock) GetBlocks(start, end uint64) ([]*types.Block, error) {
	return m.Blocks[start:end], nil
}

func (m *lightLinkMock) SimulateAddBlock(block *types.Block) {
	m.Blocks = append(m.Blocks, block)
	m.Height++
}

func (m *lightLinkMock) GetOutputV0(last *types.Header) (OutputV0, error) {
	return OutputV0{}, nil
}

func (m *lightLinkMock) GetProof(address common.Address, keys []string, height uint64) (*RawProof, error) {
	return nil, nil
}
