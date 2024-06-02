package ethereum

import (
	"fmt"
	"hummingbird/node/contracts"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Challenge interface {
	GetChallengeFee() (*big.Int, error)
	GetDataRootInclusionChallenge(block common.Hash, pointerIndex uint8, shareIndex uint32) (contracts.ChallengeDaInfo, error)
	ChallengeDataRootInclusion(index uint64, pointerIndex uint8, shareIndex uint32) (*types.Transaction, common.Hash, error)
	DefendDataRootInclusion(common.Hash, challengeContract.SharesProof) (*types.Transaction, error)
	SettleDataRootInclusion(common.Hash) (*types.Transaction, error)
	FilterChallengeDAUpdate(opts *bind.FilterOpts, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (*challengeContract.ChallengeChallengeDAUpdateIterator, error)
	DefendL2Header(common.Hash, common.Hash, common.Hash) (*types.Transaction, error)
	GetL2HeaderChallengeHash(common.Hash, *big.Int) (common.Hash, error)
	GetL2HeaderChallenge(common.Hash) (contracts.L2HeaderChallengeInfo, error)
	FilterL2HeaderChallengeUpdate(opts *bind.FilterOpts, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (*challengeContract.ChallengeL2HeaderChallengeUpdateIterator, error)
	GetChallengeWindow() (*big.Int, error)
	GetChallengeWindowBlockRanges() ([][]uint64, error)
	DataRootInclusionChallengeKey(opts *bind.CallOpts, blockHash common.Hash, pointerIndex uint8, shareIndex uint32) (common.Hash, error)
}

var _ Challenge = &Client{} // Ensure Client implements Challenge

func (c *Client) GetChallengeFee() (*big.Int, error) {
	return c.challenge.ChallengeFee(nil)
}

func (c *Client) GetDataRootInclusionChallenge(blockHash common.Hash, pointerIndex uint8, shareIndex uint32) (contracts.ChallengeDaInfo, error) {
	key, err := c.challenge.DataRootInclusionChallengeKey(nil, blockHash, pointerIndex, shareIndex)
	if err != nil {
		return contracts.ChallengeDaInfo{}, fmt.Errorf("failed to get data root inclusion challenge key: %w", err)
	}

	res, err := c.challenge.DaChallenges(nil, key)
	if err != nil {
		return contracts.ChallengeDaInfo{}, fmt.Errorf("failed to get data root inclusion challenge: %w", err)
	}

	return contracts.ChallengeDaInfo{
		BlockIndex: res.BlockIndex,
		Challenger: res.Challenger.Hex(),
		Expiry:     res.Expiry,
		Status:     res.Status,
	}, nil
}

func (c *Client) ChallengeDataRootInclusion(index uint64, pointerIndex uint8, shareIndex uint32) (*types.Transaction, common.Hash, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to create transactor: %w", err)
	}

	// set transactions fee
	fee, err := c.GetChallengeFee()
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to get challenge fee: %w", err)
	}
	transactor.Value = fee

	// get index hash
	blockHash, err := c.canonicalStateChain.Chain(nil, big.NewInt(int64(index)))
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to get hash for block %d: %w", index, err)
	}

	tx, err := c.challenge.ChallengeDataRootInclusion(transactor, big.NewInt(int64(index)), pointerIndex, shareIndex)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to challenge data root inclusion: %w", err)
	}

	return tx, blockHash, nil
}

func (c *Client) DefendDataRootInclusion(blockHash common.Hash, proof challengeContract.SharesProof) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := c.challenge.DefendDataRootInclusion(transactor, blockHash, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to defend data root inclusion: %w", err)
	}

	return tx, nil
}

func (c *Client) SettleDataRootInclusion(blockHash common.Hash) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := c.challenge.SettleDataRootInclusion(transactor, blockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to settle data root inclusion: %w", err)
	}

	return tx, nil
}

func (c *Client) FilterChallengeDAUpdate(opts *bind.FilterOpts, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (*challengeContract.ChallengeChallengeDAUpdateIterator, error) {
	return c.challenge.FilterChallengeDAUpdate(opts, _blockHash, _blockIndex, _status)
}

func (c *Client) DefendL2Header(blockHash, rootHash, headerHash common.Hash) (*types.Transaction, error) {
	transactor, err := c.transactor()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return c.challenge.DefendL2Header(transactor, blockHash, rootHash, headerHash)
}

func (c *Client) GetL2HeaderChallengeHash(rblockHash common.Hash, l2Num *big.Int) (common.Hash, error) {
	return c.challenge.L2HeaderChallengeHash(nil, rblockHash, l2Num)
}

func (c *Client) GetL2HeaderChallenge(challengeHash common.Hash) (contracts.L2HeaderChallengeInfo, error) {
	res, err := c.challenge.L2HeaderChallenges(nil, challengeHash)
	if err != nil {
		return contracts.L2HeaderChallengeInfo{}, fmt.Errorf("failed to get L2 header challenge: %w", err)
	}

	return contracts.L2HeaderChallengeInfo{
		Header:       res.Header,
		PrevHeader:   res.PrevHeader,
		ChallengeEnd: res.ChallengeEnd,
		Challenger:   res.Challenger,
		Status:       res.Status,
	}, nil
}

func (c *Client) FilterL2HeaderChallengeUpdate(opts *bind.FilterOpts, _blockHash [][32]byte, _blockIndex []*big.Int, _status []uint8) (*challengeContract.ChallengeL2HeaderChallengeUpdateIterator, error) {
	return c.challenge.FilterL2HeaderChallengeUpdate(opts, _blockHash, _blockIndex, _status)
}

func (c *Client) GetChallengeWindow() (*big.Int, error) {
	return c.challenge.ChallengeWindow(nil)
}

func (c *Client) DataRootInclusionChallengeKey(opts *bind.CallOpts, blockHash common.Hash, pointerIndex uint8, shareIndex uint32) (common.Hash, error) {
	key, err := c.challenge.DataRootInclusionChallengeKey(opts, blockHash, pointerIndex, shareIndex)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get data root inclusion challenge key: %w", err)
	}
	return common.BytesToHash(key[:]), nil
}

// Returns the block range required to log scan for open challenges.
// Useful for scanning logs for pending challenges due to eth_getLogs
// range limitations. Ranges are split into 10k block chunks to avoid
// hitting the eth_getLogs limit.
func (c *Client) GetChallengeWindowBlockRanges() ([][]uint64, error) {
	window, err := c.GetChallengeWindow() // seconds
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge window: %w", err)
	}

	// divide window by the optimistic average block time
	// to find the number of L1 blocks we need to scan
	windowsMs := window.Mul(window, big.NewInt(1000))
	numBlocksToScan := window.Div(windowsMs, big.NewInt(int64(c.opts.BlockTime)))

	// get the current block number
	currentBlock, err := c.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}

	// subtract the number of blocks we need to scan from the current block
	// to find the block where the challenge window has closed
	startBlock := currentBlock - numBlocksToScan.Uint64()

	// fill array with ranges of blocks to scan
	var blockRanges [][]uint64

	blockSize := uint64(10000)
	for startBlock+blockSize < currentBlock {
		blockRanges = append(blockRanges, []uint64{startBlock, startBlock + blockSize})
		startBlock += blockSize + 1
	}
	blockRanges = append(blockRanges, []uint64{startBlock, currentBlock})

	return blockRanges, nil
}
