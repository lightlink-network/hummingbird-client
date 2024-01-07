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

type blockResult struct {
	*types.Header `json:",inline"`
	Transactions  types.Transactions `json:"transactions"`
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

func bindJsonTx(from any, target *types.Transaction) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return target.UnmarshalJSON(b)
}
