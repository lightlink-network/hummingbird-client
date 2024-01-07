package node

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LinkLink is a client for the LightLink layer 2 network.
type LightLink interface {
	GetHeight() (uint64, error) // GetHeight returns the current height of the lightlink network.
	GetBlock(height uint64) (*types.Block, error)
	GetBlocks(start, end uint64) ([]*types.Block, error)
}

type LightLinkRPC struct {
	client *ethclient.Client
	opts   LightLinkRPCOpts
}

type LightLinkRPCOpts struct {
	Endpoint string
	Delay    time.Duration
	Logger   *slog.Logger
}

// NewLightLinkRPC returns a new LightLinkRPC client.
func NewLightLinkRPC(opts LightLinkRPCOpts) (*LightLinkRPC, error) {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	log := opts.Logger.With("func", "NewLightLinkRPC")

	client, err := ethclient.Dial(opts.Endpoint)
	if err != nil {
		log.Error("Failed to connect to LightLink", "error", err)
		return nil, err
	}

	return &LightLinkRPC{client: client, opts: opts}, nil
}

func (l *LightLinkRPC) GetHeight() (uint64, error) {
	return l.client.BlockNumber(context.Background())
}

func (l *LightLinkRPC) GetBlock(height uint64) (*types.Block, error) {
	return l.client.BlockByNumber(context.Background(), new(big.Int).SetUint64(height))
}

func (l *LightLinkRPC) GetBlocks(start, end uint64) ([]*types.Block, error) {
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
