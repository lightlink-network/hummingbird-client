package rollup

import (
	"encoding/json"
	"errors"
	"fmt"
	"hummingbird/node"
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common"
)

type RollupInfo struct {
	RollupHeight     uint64 `pretty:"Rollup Height"`
	L2BlocksRolledUp uint64 `pretty:"L2 Blocks Rolled Up"`
	L2BlocksTodo     uint64 `pretty:"L2 Blocks Todo"`

	LatestRollup struct {
		Hash                                 common.Hash `pretty:"Hash"`
		BundleSize                           uint64      `pretty:"Bundle Size"`
		*contracts.CanonicalStateChainHeader `pretty:"Header"`
	} `pretty:"Latest Rollup Block"`

	DataAvailability struct {
		CelestiaHeight   uint64   `pretty:"Celestia Height"`
		CelestiaDataRoot [32]byte `pretty:"Celestia Data Root"`
		CelestiaTx       string   `pretty:"Celestia Tx"`
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
	info.DataAvailability.CelestiaDataRoot = latestRollupHead.CelestiaDataRoot
	info.DataAvailability.CelestiaTx = "Unknown"

	pointer, err := r.getDAPointer(latestRollupHash)
	if err != nil || pointer == nil {
		r.Opts.Logger.Warn("Failed to get celestia pointer", "error", err)
	} else {
		info.DataAvailability.CelestiaTx = pointer.TxHash.Hex()
	}

	return info, nil
}

func (r *Rollup) getDAPointer(hash common.Hash) (*node.CelestiaPointer, error) {
	if r.Store == nil {
		return nil, errors.New("no store")
	}

	key := append([]byte("pointer_"), hash[:]...)
	buf, err := r.Store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get celestia pointer from store: %w", err)
	}

	pointer := &node.CelestiaPointer{}
	err = json.Unmarshal(buf, pointer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal celestia pointer: %w", err)
	}

	return pointer, nil
}

type RollupBlockInfo struct {
	Hash                                 common.Hash `pretty:"Hash"`
	BundleSize                           uint64      `pretty:"Bundle Size"`
	*contracts.CanonicalStateChainHeader `pretty:"Header"`

	DataAvailability struct {
		CelestiaHeight   uint64   `pretty:"Celestia Height"`
		CelestiaDataRoot [32]byte `pretty:"Celestia Data Root"`
		CelestiaTx       string   `pretty:"Celestia Tx"`
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
	rbi.DataAvailability.CelestiaDataRoot = header.CelestiaDataRoot
	rbi.DataAvailability.CelestiaTx = "Unknown"

	// get celestia pointer
	pointer, err := r.getDAPointer(hash)
	if err != nil || pointer == nil {
		r.Opts.Logger.Warn("Failed to get celestia pointer", "error", err)
	} else {
		rbi.DataAvailability.CelestiaTx = pointer.TxHash.Hex()
	}

	return rbi, nil
}
