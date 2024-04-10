package ethereum

import (
	"fmt"
	"math/big"

	chainOracleContract "hummingbird/node/contracts/ChainOracle.sol"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ChainOracle interface {
	ProvideShares(rblock common.Hash, pointerIndex uint8, shareProof *chainOracleContract.SharesProof, attProof chainOracleContract.AttestationProof) (*types.Transaction, error)
	ProvideHeader(rblock common.Hash, shareData [][]byte, ranges []chainOracleContract.ChainOracleShareRange) (*types.Transaction, error)
	ProvideLegacyTx(rblock common.Hash, shareData [][]byte, ranges []chainOracleContract.ChainOracleShareRange) (*types.Transaction, error)
	AlreadyProvidedShares(rblock common.Hash, shareData [][]byte) (bool, error)
	AlreadyProvidedHeader(l2Hash common.Hash) (bool, error)
}

func (c *Client) ProvideShares(rblock common.Hash, pointerIndex uint8, shareProof *chainOracleContract.SharesProof, attestationProof chainOracleContract.AttestationProof) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return c.chainLoader.ProvideShares(transactor, rblock, pointerIndex, *shareProof)
}

func (c *Client) ProvideHeader(rblock common.Hash, shareData [][]byte, ranges []chainOracleContract.ChainOracleShareRange) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	sharekey, err := c.chainLoader.ShareKey(nil, rblock, shareData)
	if err != nil {
		return nil, fmt.Errorf("failed to get share key: %w", err)
	}

	// check shares are found
	s, err := c.chainLoader.Shares(nil, sharekey, big.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("failed checking shares were deployed: %w", err)
	}

	if len(s) == 0 {
		return nil, fmt.Errorf("failed checking shares: shares not found")
	}

	return c.chainLoader.ProvideHeader(transactor, sharekey, ranges)
}

func (c *Client) ProvideLegacyTx(rblock common.Hash, shareData [][]byte, ranges []chainOracleContract.ChainOracleShareRange) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	sharekey, err := c.chainLoader.ShareKey(nil, rblock, shareData)
	if err != nil {
		return nil, fmt.Errorf("failed to get share key: %w", err)
	}

	// check shares are found
	s, err := c.chainLoader.Shares(nil, sharekey, big.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("failed checking shares were deployed: %w", err)
	}

	if len(s) == 0 {
		return nil, fmt.Errorf("failed checking shares: shares not found")
	}

	return c.chainLoader.ProvideLegacyTx(transactor, sharekey, ranges)
}

func (c *Client) AlreadyProvidedShares(rblock common.Hash, shareData [][]byte) (bool, error) {
	sharekey, err := c.chainLoader.ShareKey(nil, rblock, shareData)
	if err != nil {
		return false, fmt.Errorf("failed to get share key: %w", err)
	}

	// check shares are found
	s, err := c.chainLoader.Shares(nil, sharekey, big.NewInt(0))
	if err != nil {
		return false, fmt.Errorf("failed checking shares were deployed: %w", err)
	}

	return len(s) > 0, nil
}

func (c *Client) AlreadyProvidedHeader(l2Hash common.Hash) (bool, error) {
	h, err := c.chainLoader.GetHeader(nil, l2Hash)
	if err != nil {
		return false, fmt.Errorf("failed to get header: %w", err)
	}

	return h.Number.Uint64() > 0, nil
}
