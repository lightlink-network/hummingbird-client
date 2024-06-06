package ethereum

import (
	"context"
	"fmt"
	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type CanonicalStateChain interface {
	GetRollupHeight() (uint64, error)                                                                         // Get the current rollup block height.
	GetHeight() (uint64, error)                                                                               // Get the current block height of the Ethereum network.
	GetRollupHead() (canonicalStateChainContract.CanonicalStateChainHeader, error)                            // Get the latest rollup block header in the CanonicalStateChain.sol contract.
	PushRollupHead(header *canonicalStateChainContract.CanonicalStateChainHeader) (*types.Transaction, error) // Push a new rollup block header to the CanonicalStateChain.sol contract.
	GetRollupHeader(index uint64) (canonicalStateChainContract.CanonicalStateChainHeader, error)              // Get the rollup block header at the given index from the CanonicalStateChain.sol contract.
	GetRollupHeaderByHash(hash common.Hash) (canonicalStateChainContract.CanonicalStateChainHeader, error)    // Get the rollup block header with the given hash from the CanonicalStateChain.sol contract.
	Wait(txHash common.Hash) (*types.Receipt, error)                                                          // Wait for a transaction to be mined.
	GetPublisher() (common.Address, error)                                                                    // Get the address of the publisher of the CanonicalStateChain.sol contract.
	HashHeader(header *canonicalStateChainContract.CanonicalStateChainHeader) (common.Hash, error)            // Hash a rollup block header.
}

// GetRollupHeight returns the current rollup block height.
func (c *Client) GetRollupHeight() (uint64, error) {
	h, err := c.canonicalStateChain.ChainHead(nil)
	if err != nil {
		return 0, err
	}

	return h.Uint64(), nil
}

func (c *Client) GetHeight() (uint64, error) {
	return c.client.BlockNumber(context.Background())
}

// GetRollupHead returns the latest rollup block header.
func (c *Client) GetRollupHead() (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return c.canonicalStateChain.GetHead(nil)
}

// PushRollupHead pushes a new rollup block header.
func (c *Client) PushRollupHead(header *canonicalStateChainContract.CanonicalStateChainHeader) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return c.canonicalStateChain.PushBlock(transactor, *header)
}

// GetRollupHeader returns the rollup block header at the given index.
func (c *Client) GetRollupHeader(index uint64) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return c.canonicalStateChain.GetHeaderByNum(nil, big.NewInt(int64(index)))
}

// GetRollupHeaderByHash returns the rollup block header with the given hash.
func (c *Client) GetRollupHeaderByHash(hash common.Hash) (canonicalStateChainContract.CanonicalStateChainHeader, error) {
	return c.canonicalStateChain.GetHeaderByHash(nil, hash)
}

// Wait waits for a transaction to be mined and returns the receipt
func (c *Client) Wait(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()

	for {
		c.logger.Debug("Waiting for transaction to be mined...", "txHash", txHash.Hex())
		// try to get the receipt
		receipt, err := c.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			// if the receipt is not found, keep checking
			if err.Error() == "not found" {
				time.Sleep(10 * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
		c.logger.Debug("Transaction mined successfully!", "txHash", txHash.Hex())
		return receipt, nil
	}
}

func (c *Client) GetPublisher() (common.Address, error) {
	return c.canonicalStateChain.Publisher(nil)
}

func (c *Client) HashHeader(header *canonicalStateChainContract.CanonicalStateChainHeader) (common.Hash, error) {
	return c.canonicalStateChain.CalculateHeaderHash(nil, *header)
}
