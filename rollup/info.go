package rollup

import (
	"fmt"
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common"

	canonicalStateChainContract "hummingbird/node/contracts/CanonicalStateChain.sol"
)

type RollupInfo struct {
	RollupHeight     uint64 `pretty:"Rollup Height"`
	L2BlocksRolledUp uint64 `pretty:"L2 Blocks Rolled Up"`
	L2BlocksTodo     uint64 `pretty:"L2 Blocks Todo"`

	LatestRollup struct {
		Hash                                                   common.Hash `pretty:"Hash"`
		BundleSize                                             uint64      `pretty:"Bundle Size"`
		*canonicalStateChainContract.CanonicalStateChainHeader `pretty:"Header"`
	} `pretty:"Latest Rollup Block"`

	DataAvailability struct {
		CelestiaHeight     uint64   `pretty:"Celestia Height"`
		CelestiaShareStart uint64   `pretty:"Shares Start"`
		CelestiaShareLen   uint64   `pretty:"Shares"`
		CelestiaTx         [32]byte `pretty:"Celestia Tx"`
	} `pretty:"Data Availability"`
}

func (r *Rollup) GetInfo() (*RollupInfo, error) {
	info := &RollupInfo{}

	// get rollup height
	rollupHeight, err := r.Ethereum.GetRollupHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get rollup height: %w", err)
	}

	// get latest rollup head
	latestRollupHead, err := r.Ethereum.GetRollupHead()
	if err != nil {
		return nil, fmt.Errorf("failed to get rollup head: %w", err)
	}

	// hash latest rollup head
	latestRollupHash, err := contracts.HashCanonicalStateChainHeader(&latestRollupHead)
	if err != nil {
		return nil, fmt.Errorf("failed to hash rollup head: %w", err)
	}

	// get previous rollup header
	prevRollupHeader, err := r.Ethereum.GetRollupHeaderByHash(latestRollupHead.PrevHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous rollup header: %w", err)
	}

	// get bundle size
	bundleSize := latestRollupHead.L2Height - prevRollupHeader.L2Height

	// get genesis header
	genesisHeader, err := r.Ethereum.GetRollupHeader(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis header: %w", err)
	}

	// get l2 blocks rolled up
	l2BlocksRolledUp := latestRollupHead.L2Height - genesisHeader.L2Height

	// get layer 2 height
	l2Height, err := r.LightLink.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get layer 2 height: %w", err)
	}

	// get l2 blocks todo
	l2BlocksTodo := l2Height - latestRollupHead.L2Height

	info.RollupHeight = rollupHeight
	info.LatestRollup.Hash = latestRollupHash
	info.LatestRollup.CanonicalStateChainHeader = &latestRollupHead
	info.LatestRollup.BundleSize = bundleSize
	info.L2BlocksRolledUp = l2BlocksRolledUp
	info.L2BlocksTodo = l2BlocksTodo

	// get data availability
	info.DataAvailability.CelestiaHeight = latestRollupHead.CelestiaHeight
	info.DataAvailability.CelestiaShareStart = latestRollupHead.CelestiaShareStart
	info.DataAvailability.CelestiaShareLen = latestRollupHead.CelestiaShareLen

	return info, nil
}

type RollupBlockInfo struct {
	Hash                                                   common.Hash `pretty:"Hash"`
	BundleSize                                             uint64      `pretty:"Bundle Size"`
	*canonicalStateChainContract.CanonicalStateChainHeader `pretty:"Header"`

	DataAvailability struct {
		CelestiaHeight     uint64 `pretty:"Celestia Height"`
		CelestiaShareStart uint64 `pretty:"Shares Start"`
		CelestiaShareLen   uint64 `pretty:"Shares"`
	} `pretty:"Data Availability"`

	Distance struct {
		FromLatestInEpochs   uint64 `pretty:"From Latest Epoch"`
		FromLatestInL1height uint64 `pretty:"From Latest L1 Height"`
		FromLatestInL2height uint64 `pretty:"From Latest L2 Height"`
	} `pretty:"Distance"`
}

func (r *Rollup) GetBlockInfo(hash common.Hash) (*RollupBlockInfo, error) {
	info, err := r.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get rollup info: %w", err)
	}

	rbi := &RollupBlockInfo{}

	// get rollup header
	header, err := r.Ethereum.GetRollupHeaderByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get rollup header: %w", err)
	}

	// get previous rollup header
	prevRollupHeader, err := r.Ethereum.GetRollupHeaderByHash(header.PrevHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous rollup header: %w", err)
	}

	// hash rollup header
	hash, err = contracts.HashCanonicalStateChainHeader(&header)
	if err != nil {
		return nil, fmt.Errorf("failed to hash rollup header: %w", err)
	}

	// get l2 height
	llHeight, err := r.LightLink.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get l2 height: %w", err)
	}

	// get l1 height
	ethHeight, err := r.Ethereum.GetHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get l1 height: %w", err)
	}

	// set distance from latest
	rbi.Distance.FromLatestInEpochs = info.LatestRollup.Epoch - header.Epoch
	rbi.Distance.FromLatestInL1height = ethHeight - header.Epoch
	rbi.Distance.FromLatestInL2height = llHeight - header.L2Height

	// set header
	rbi.Hash = hash
	rbi.BundleSize = header.L2Height - prevRollupHeader.L2Height
	rbi.CanonicalStateChainHeader = &header

	// set data availability
	rbi.DataAvailability.CelestiaHeight = header.CelestiaHeight
	rbi.DataAvailability.CelestiaShareStart = header.CelestiaShareStart
	rbi.DataAvailability.CelestiaShareLen = header.CelestiaShareLen

	return rbi, nil
}
