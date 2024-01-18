package node

import (
	"encoding/json"
	"fmt"
	"hummingbird/node/jsonrpc"
	"hummingbird/node/lightlink"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// LinkLink is a client for the LightLink layer 2 network.
type LightLink interface {
	GetHeight() (uint64, error) // GetHeight returns the current height of the lightlink network.
	GetBlock(height uint64) (*types.Block, error)
	GetBlocks(start, end uint64) ([]*types.Block, error)
}

type LightLinkClientOpts struct {
	Endpoint string
	Delay    time.Duration
	Logger   *slog.Logger
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

	result := resp.Result.(map[string]interface{})
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
	for i := start; i < end; i++ {
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
