package node

import (
	"encoding/json"
	"fmt"
	"hummingbird/node/jsonrpc"
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
	log := opts.Logger.With("func", "NewLightLinkClient")

	client, err := jsonrpc.NewClient(opts.Endpoint)
	if err != nil {
		log.Error("Failed to connect to LightLink", "error", err)
		return nil, err
	}

	return &LightLinkClient{client: client, opts: opts}, nil
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
		tx := &types.Transaction{}
		err := bindJsonTx(v, tx)
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
	log := slog.Default().With("func", "GetBlocks")

	var blocks []*types.Block
	for i := start; i <= end; i++ {
		block, err := l.GetBlock(i)
		if err != nil {
			log.Error("Failed to get block", "height", i, "error", err)
			return nil, err
		}
		blocks = append(blocks, block)

		// delay between requests
		if l.opts.Delay > 0 {
			time.Sleep(l.opts.Delay)
		}
	}

	return blocks, nil
}

func bindJsonTx(from any, target *types.Transaction) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return target.UnmarshalJSON(b)
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
