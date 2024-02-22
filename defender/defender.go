package defender

import (
	"context"
	"fmt"
	"hummingbird/node"
	"strings"
	"time"

	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"hummingbird/node/contracts"
	challengeContract "hummingbird/node/contracts/Challenge.sol"
)

type Opts struct {
	Logger      *slog.Logger
	WorkerDelay time.Duration
}

type Defender struct {
	*node.Node
	Opts *Opts
}

func NewDefender(node *node.Node, opts *Opts) *Defender {
	return &Defender{Node: node, Opts: opts}
}

// Start starts the defender.
//
// It will:
//  1. Start a goroutine to Scan Challenge.sol historic events, find any pending
//     DA challenges that were missed and defend them.
//  2. In main thread, watch Challenge.sol for new DA challenges and defend them.
func (d *Defender) Start() error {
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		errChan <- d.startHistoricWorker(ctx)
	}()

	go func() {
		errChan <- d.watchAndDefendDAChallenges(ctx)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (d *Defender) startHistoricWorker(ctx context.Context) error {
	ticker := time.NewTicker(d.Opts.WorkerDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := d.scanAndDefendHistoricChallenges(); err != nil {
				return fmt.Errorf("error retrying historic DA challenges: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Watches the Challenge.sol contract for new DA challenges and defends them.
func (d *Defender) watchAndDefendDAChallenges(ctx context.Context) error {
	challenges := make(chan *challengeContract.ChallengeChallengeDAUpdate)
	errChan := make(chan error, 1)

	subscription, err := d.Ethereum.WatchChallengesDA(challenges)
	if err != nil {
		return fmt.Errorf("error starting WatchChallengesDA: %w", err)
	}
	defer subscription.Unsubscribe()

	d.Opts.Logger.Info("Defender is watching for DA challenges")

	go func() {
		<-ctx.Done()
		close(challenges)
	}()

	var wg sync.WaitGroup

	for challenge := range challenges {
		wg.Add(1)
		go func(challenge *challengeContract.ChallengeChallengeDAUpdate) {
			defer wg.Done()

			if err := d.handleDAChallenge(challenge); err != nil {
				errChan <- err
				return
			}
		}(challenge)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Handles a DA challenge by attempting to defend it.
//
// If the challenged data root is not yet available, it will be ignored
// and retried later by the scanAndDefendHistoricChallenges() worker function.
func (d *Defender) handleDAChallenge(challenge *challengeContract.ChallengeChallengeDAUpdate) error {
	// we are only interested in challenges that have been initiated by a challenger, ready to be defended
	if challenge.Status != 1 {
		return nil
	}

	blockHash := common.BytesToHash(challenge.BlockHash[:])
	statusString := contracts.StatusString(challenge.Status)

	log := d.Opts.Logger.With(
		"blockHash", blockHash.Hex(),
		"blockIndex", challenge.BlockIndex,
		"expiry", time.Unix(challenge.Expiry.Int64(), 0).Format(time.RFC1123Z),
		"statusEnum", challenge.Status,
		"statusString", statusString,
	)
	log.Info("Pending DA challenge log event found")

	// attempt to defend the challenge by submitting a tx to the Challenge contract
	tx, err := d.DefendDA(challenge.BlockHash)
	if err != nil {
		if strings.Contains(err.Error(), "no data commitment has been generated for the provided height") {
			log.Info("Pending DA challenge is awaiting data commitment from Celestia validators, will retry later")
			return nil
		} else {
			return fmt.Errorf("error defending DA challenge: %w", err)
		}
	}

	log.Info("Pending DA challenge defended successfully", "tx", tx.Hash().Hex())
	return nil
}

// Attempts to defend a DA challenge for the given block hash.
//
// Queries Celestia for a proof of data availability and submits a tx to the Challenge contract.
func (d *Defender) DefendDA(block common.Hash) (*types.Transaction, error) {
	proof, err := d.GetDAProof(block)
	if err != nil {
		return nil, fmt.Errorf("failed to prove data availability: %w", err)
	}
	return d.Ethereum.DefendDataRootInclusion(block, proof)
}

// Gets the Celestia pointer for the given block hash and queries Celestia for a proof
// of data availability.
func (d *Defender) GetDAProof(block common.Hash) (*node.CelestiaProof, error) {
	pointer, err := d.GetDAPointer(block)
	if err != nil {
		return nil, fmt.Errorf("failed to get Celestia pointer: %w", err)
	}
	if pointer == nil {
		return nil, fmt.Errorf("no Celestia pointer found")
	}
	return d.Celestia.GetProof(pointer)
}

// Scans the Challenge.sol contract for historic DA challenges and attempts to defend them.
//
// This function runs every opts.WorkerDelay and will scan all historic Challenge.sol challenge
// logs. This ensure we don't miss any challenges that were initiated when offline. It also
// allows defenders to retry challenges that failed to be defended i.e. due to data commitments
// not being available yet.
func (d *Defender) scanAndDefendHistoricChallenges() error {
	d.Opts.Logger.Debug("Starting log scan for historic pending DA challenges")

	h, err := d.Ethereum.GetRollupHeader(uint64(1))
	if err != nil {
		return fmt.Errorf("error getting rollup header: %w", err)
	}

	opts := &bind.FilterOpts{
		Start: h.Epoch,
	}

	challenges, err := d.Ethereum.FilterChallengeDAUpdate(opts, nil, nil, []uint8{1})
	if err != nil {
		return fmt.Errorf("error filtering challenges: %w", err)
	}

	// iterate through historic challenges events
	for challenges.Next() {
		challenge := challenges.Event

		// check if challenge has already been defended by checking the current status
		challengeInfo, err := d.Ethereum.GetDataRootInclusionChallenge(challenge.BlockHash)
		if err != nil {
			d.Opts.Logger.Error("error getting data root inclusion challenge:", "error", err)
			continue
		}

		if challengeInfo.Status != 1 {
			continue
		}

		err = d.handleDAChallenge(challenge)
		if err != nil {
			continue
		}

	}
	d.Opts.Logger.Debug("Finished log scan for historic pending DA challenges")

	return nil
}

func (d *Defender) ProvideL2Header(rblock common.Hash, l2Block common.Hash, skipShares bool) (*types.Transaction, error) {

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundle, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Get a pointer to the shares that contain the header
	sharePointer, err := bundle.FindHeaderShares(l2Block, d.Namespace())
	if err != nil {
		return nil, fmt.Errorf("error finding header shares in the bundle: %w", err)
	}

	// Get proof the shares are in the bundle
	shareProof, err := d.Celestia.GetSharesProof(&node.CelestiaPointer{
		Height:     rheader.CelestiaHeight,
		ShareStart: rheader.CelestiaShareStart,
		ShareLen:   rheader.CelestiaShareLen,
	}, sharePointer)
	if err != nil {
		return nil, fmt.Errorf("error getting share proof: %w", err)
	}

	// Get proof the data is available
	celProof, err := d.GetDAProof(rblock)
	if err != nil {
		return nil, fmt.Errorf("error proving data availability: %w", err)
	}

	// Provide the shares
	if !skipShares {
		tx, err := d.Ethereum.ProvideShares(rblock, shareProof, celProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())
	}

	// Finally, provide the header
	return d.Ethereum.ProvideHeader(rblock, shareProof.Data, *sharePointer)
}

func (d *Defender) ProvideL2Tx(rblock common.Hash, l2Tx common.Hash, skipShares bool) (*types.Transaction, error) {

	// Download the rollup block and bundle from L1 and
	// Celestia
	rheader, bundle, err := d.Node.FetchRollupBlock(rblock)
	if err != nil {
		return nil, fmt.Errorf("error fetching rollup block: %w", err)
	}

	// Get a pointer to the shares that contain the transaction
	sharePointer, err := bundle.FindTxShares(l2Tx, d.Namespace())
	if err != nil {
		return nil, fmt.Errorf("error finding tx shares in the bundle: %w", err)
	}

	// Get proof the shares are in the bundle
	shareProof, err := d.Celestia.GetSharesProof(&node.CelestiaPointer{
		Height:     rheader.CelestiaHeight,
		ShareStart: rheader.CelestiaShareStart,
		ShareLen:   rheader.CelestiaShareLen,
	}, sharePointer)
	if err != nil {
		return nil, fmt.Errorf("error getting share proof: %w", err)
	}

	// Get proof the data is available
	celProof, err := d.GetDAProof(rblock)
	if err != nil {
		return nil, fmt.Errorf("error proving data availability: %w", err)
	}

	// Provide the shares
	if !skipShares {
		tx, err := d.Ethereum.ProvideShares(rblock, shareProof, celProof)
		if err != nil {
			return nil, fmt.Errorf("error providing shares: %w", err)
		}
		d.Opts.Logger.Info("Provided shares", "tx", tx.Hash().Hex(), "block", rblock.Hex(), "shares", len(shareProof.Data))
		d.Ethereum.Wait(tx.Hash())
	}

	// Finally, provide the transaction
	return d.Ethereum.ProvideLegacyTx(rblock, shareProof.Data, *sharePointer)
}
