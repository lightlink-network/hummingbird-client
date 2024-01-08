package rollup

import (
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common"
)

type RollupInfo struct {
	RollupHeight uint64
	LatestRollup struct {
		Hash common.Hash
		*contracts.CanonicalStateChainHeader
		BundleSize uint64
	}
	L2BlocksRolledUp uint64
	L2BlocksTodo     uint64
}

func (r *Rollup) GetInfo() (*RollupInfo, error) {
	log := r.Opts.Logger.With("func", "GetInfo")
	info := &RollupInfo{}

	// get rollup height
	rollupHeight, err := r.Ethereum.GetRollupHeight()
	if err != nil {
		log.Error("Failed to get rollup height", "error", err)
		return nil, err
	}

	// get latest rollup head
	latestRollupHead, err := r.Ethereum.GetRollupHead()
	if err != nil {
		log.Error("Failed to get rollup head", "error", err)
		return nil, err
	}

	// hash latest rollup head
	latestRollupHash, err := contracts.HashCanonicalStateChainHeader(&latestRollupHead)
	if err != nil {
		log.Error("Failed to hash rollup head", "error", err)
		return nil, err
	}

	// get previous rollup header
	prevRollupHeader, err := r.Ethereum.GetRollupHeaderByHash(latestRollupHead.PrevHash)
	if err != nil {
		log.Error("Failed to get previous rollup header", "error", err)
		return nil, err
	}

	// get bundle size
	bundleSize := latestRollupHead.L2Height - prevRollupHeader.L2Height

	// get genesis header
	genesisHeader, err := r.Ethereum.GetRollupHeader(0)
	if err != nil {
		log.Error("Failed to get genesis header", "error", err)
		return nil, err
	}

	// get l2 blocks rolled up
	l2BlocksRolledUp := latestRollupHead.L2Height - genesisHeader.L2Height

	// get layer 2 height
	l2Height, err := r.LightLink.GetHeight()
	if err != nil {
		log.Error("Failed to get layer 2 height", "error", err)
		return nil, err
	}

	// get l2 blocks todo
	l2BlocksTodo := l2Height - latestRollupHead.L2Height

	info.RollupHeight = rollupHeight
	info.LatestRollup.Hash = latestRollupHash
	info.LatestRollup.CanonicalStateChainHeader = &latestRollupHead
	info.LatestRollup.BundleSize = bundleSize
	info.L2BlocksRolledUp = l2BlocksRolledUp
	info.L2BlocksTodo = l2BlocksTodo

	return info, nil
}
