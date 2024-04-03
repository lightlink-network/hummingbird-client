package ethereum

import (
	"context"
	"fmt"
	"math/big"

	blobstreamXContract "hummingbird/node/contracts/BlobstreamX.sol"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type BlobstreamX interface {
	FilterDataCommitmentStored(opts *bind.FilterOpts, startBlock []uint64, endBlock []uint64, dataCommitment [][32]byte) (*blobstreamXContract.BlobstreamXDataCommitmentStoredIterator, error)
	DAVerify(proofNonce *big.Int, tuple blobstreamXContract.DataRootTuple, proof blobstreamXContract.BinaryMerkleProof) (bool, error)
	GetBlobstreamCommitment(height int64) (*blobstreamXContract.BlobstreamXDataCommitmentStored, error)
}

func (c *Client) FilterDataCommitmentStored(opts *bind.FilterOpts, startBlock []uint64, endBlock []uint64, dataCommitment [][32]byte) (*blobstreamXContract.BlobstreamXDataCommitmentStoredIterator, error) {
	return c.blobstreamX.FilterDataCommitmentStored(opts, startBlock, endBlock, dataCommitment)
}

func (c *Client) DAVerify(proofNonce *big.Int, tuple blobstreamXContract.DataRootTuple, proof blobstreamXContract.BinaryMerkleProof) (bool, error) {
	return c.blobstreamX.VerifyAttestation(nil, proofNonce, tuple, proof)
}

// GetBlobstreamCommitment returns the commitment for the given celestia height.
// see https://docs.celestia.org/developers/blobstream-proof-queries
func (c *Client) GetBlobstreamCommitment(height int64) (*blobstreamXContract.BlobstreamXDataCommitmentStored, error) {
	scanRanges := c.GetChallengeWindowBlockRanges()

	lastCommitHeight := uint64(0)
	for i := 0; i < len(scanRanges); i++ {
		if len(scanRanges[i]) != 2 {
			return nil, fmt.Errorf("invalid block range")
		}

		// get all events
		events, err := c.blobstreamX.FilterDataCommitmentStored(&bind.FilterOpts{
			Context: context.Background(),
			Start:   scanRanges[i][0],
			End:     &scanRanges[i][1],
		}, nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to filter events: %w", err)
		}
		defer events.Close()

		for events.Next() {
			e := events.Event
			if e.EndBlock > lastCommitHeight {
				lastCommitHeight = e.EndBlock
			}

			if int64(e.StartBlock) <= height && height < int64(e.EndBlock) {
				return e, nil
			}
		}
		if err := events.Error(); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("no commitment found for height %d (last commitment is for %d)", height, lastCommitHeight)
}
