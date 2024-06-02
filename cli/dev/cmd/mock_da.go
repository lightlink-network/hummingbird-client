package cmd

import (
	"hummingbird/defender"
	canonicalstatechain "hummingbird/node/contracts/CanonicalStateChain.sol"
	challenge "hummingbird/node/contracts/Challenge.sol"
	"hummingbird/rollup"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

type MockDaData struct {
	RollupHash   common.Hash                                   `json:"rollupHash"`
	RollupHeader canonicalstatechain.CanonicalStateChainHeader `json:"rollupHeader"`
	DaProof      CHallengeDaData                               `json:"daProofs"`
}

type CHallengeDaData struct {
	Key                    common.Hash                  `json:"key"`
	PointerIndex           uint8                        `json:"pointerIndex"`
	ShareIndex             uint32                       `json:"shareIndex"`
	ShareProof             *challenge.SharesProof       `json:"shareProof"`
	ShareToRBlockRootProof *challenge.BinaryMerkleProof `json:"shareToRBlockRootProof"`
}

var MockDaCmd = &cobra.Command{
	Use:   "mock-da [rblock]:[pointer]:[shareIndex]",
	Short: "mock-da will output mock data for testing using real blocks",
	Long:  "mock-da will output mock data for testing using real blocks. The first argument is in the form rblock_hash:pointer_index.",
	ArgAliases: []string{
		"rblock:pointer",
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// 0. parse args
		rblockPointer := strings.Split(args[0], ":")
		rblockHash := common.HexToHash(rblockPointer[0])
		pointerIndex, _ := strconv.Atoi(rblockPointer[1])
		shareIndex, _ := strconv.Atoi(rblockPointer[2])

		// 1. make node
		n, log, err := makeNode()
		panicErr(err, "failed to create node")
		r := rollup.NewRollup(n, &rollup.Opts{
			Logger: log.With("ctx", "Rollup"),
		})

		// 2. Get rblock
		rblock, err := r.GetBlockByHash(rblockHash)
		panicErr(err, "failed to get rollup block")
		log.Info("Got rblock", "hash", rblockHash.String(), "bundles", len(rblock.Bundles))

		// 3. Create defender and get proofs
		d := defender.NewDefender(n, &defender.Opts{
			Logger: log.With("ctx", "Defender"),
		})

		key, shareProof, shareToRblockRootProof, err := d.GetDaProof(rblockHash, uint8(pointerIndex), uint32(shareIndex))
		panicErr(err, "failed to get da proof")

		// 4. Output the mock data
		out := &MockDaData{
			RollupHash:   rblockHash,
			RollupHeader: *rblock.CanonicalStateChainHeader,
			DaProof: CHallengeDaData{
				Key:                    *key,
				PointerIndex:           uint8(pointerIndex),
				ShareIndex:             uint32(shareIndex),
				ShareProof:             shareProof,
				ShareToRBlockRootProof: shareToRblockRootProof,
			},
		}

		printJSON(out)
	},
}
