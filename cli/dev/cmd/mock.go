package cmd

import (
	"hummingbird/node"
	"hummingbird/node/contracts"
	canonicalstatechain "hummingbird/node/contracts/CanonicalStateChain.sol"
	chainoracle "hummingbird/node/contracts/ChainOracle.sol"
	chainoracleContract "hummingbird/node/contracts/ChainOracle.sol"
	"hummingbird/rollup"
	"hummingbird/utils"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

type MockData struct {
	RollupHash   common.Hash                                   `json:"rollupHash"`
	RollupHeader canonicalstatechain.CanonicalStateChainHeader `json:"rollupHeader"`
	Headers      []HeaderData                                  `json:"headers"`
}

type HeaderData struct {
	Header        *types.Header                       `json:"header"`
	HeaderHash    common.Hash                         `json:"headerHash"`
	ShareProofs   chainoracle.SharesProof             `json:"shareProofs"`
	ShareRanges   []chainoracle.ChainOracleShareRange `json:"shareRanges"`
	PointerProofs []chainoracle.BinaryMerkleProof     `json:"pointerProofs"`
	Shares        [][]byte                            `json:"shares"`
}

func init() {
	RootCmd.AddCommand(MockCmd)
}

var MockCmd = &cobra.Command{
	Use:   "mock [rblock] [num]",
	Short: "mock will output mock data for testing using real blocks",
	Long:  "mock will output mock data for testing using real blocks for a given rblock hash. `num` is the number of headers proofs to generate.",
	Args:  cobra.MinimumNArgs(2),
	ArgAliases: []string{
		"rblock",
		"num",
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 0. parse args
		rblockHash := common.HexToHash(args[0])
		proofsNum, _ := strconv.Atoi(args[1])

		// 1. make node
		n, log, err := makeNode()
		panicErr(err, "failed to create node")
		r := rollup.NewRollup(n, &rollup.Opts{
			Logger: log.With("ctx", "Rollup"),
		})

		// 2. Get rblock and celestia pointer
		rblock, err := r.GetBlockByHash(rblockHash)
		panicErr(err, "failed to get rollup block")
		log.Info("Got rblock", "hash", rblockHash.String(), "bundles", len(rblock.Bundles))

		// 3. Fetch the items
		bundles := rblock.Bundles
		hds := make([]HeaderData, 0)
		log.Info("Generating mock data", "num", proofsNum)
		for i := 0; i < proofsNum; i++ {
			// select a random bundle from the rblock
			pointerIndex := rand.Intn(len(bundles))
			bundle := bundles[pointerIndex]

			// select a random header from the bundle
			header := bundle.Blocks[rand.Intn(len(bundle.Blocks))].Header()
			headerHash := utils.HashHeaderWithoutExtraData(header)

			sharePointer, err := bundle.FindHeaderShares(headerHash, r.Namespace())
			panicErr(err, "failed to find header shares")

			shareProof, err := r.Celestia.GetSharesProof(rblock.GetCelestiaPointers()[pointerIndex], sharePointer)
			panicErr(err, "failed to get share proof")

			shareProofs, err := contracts.NewShareProof(shareProof, getAttestations(r.Node, rblock.GetCelestiaPointers()[pointerIndex]))
			panicErr(err, "failed to get share proofs")
			log.Info("Got share proofs", "header", headerHash.String(), "index", i)

			blockProofs := node.GetSharesProofs(sharePointer, bundles, pointerIndex, r.Namespace())

			hds = append(hds, HeaderData{
				Header:        header,
				HeaderHash:    headerHash,
				ShareProofs:   *shareProofs,
				Shares:        sharesToBytes(sharePointer.Shares()),
				ShareRanges:   formatRanges(sharePointer),
				PointerProofs: utils.ToBinaryMerkleProof(blockProofs),
			})
		}

		out := &MockData{
			RollupHash:   rblockHash,
			RollupHeader: *rblock.CanonicalStateChainHeader,
			Headers:      hds,
		}

		printJSON(out)
	},
}

func getAttestations(n *node.Node, celPointer *node.CelestiaPointer) chainoracle.AttestationProof {
	commitment, err := n.Ethereum.GetBlobstreamCommitment(int64(celPointer.Height))
	panicErr(err, "failed to get blobstream commitment")

	celProof, err := n.Celestia.GetProof(celPointer, commitment.StartBlock, commitment.EndBlock, *commitment.ProofNonce)
	panicErr(err, "failed to get celestia proof")

	return chainoracle.AttestationProof{
		TupleRootNonce: celProof.Nonce,
		Tuple: chainoracleContract.DataRootTuple{
			Height:   celProof.Tuple.Height,
			DataRoot: celProof.Tuple.DataRoot,
		},
		Proof: chainoracle.BinaryMerkleProof{
			SideNodes: celProof.WrappedProof.SideNodes,
			Key:       celProof.WrappedProof.Key,
			NumLeaves: celProof.WrappedProof.NumLeaves,
		},
	}
}

func sharesToBytes(s []shares.Share) [][]byte {
	shareBytes := [][]byte{}
	for _, share := range s {
		shareBytes = append(shareBytes, share.ToBytes())
	}
	return shareBytes
}

func formatRanges(sp *node.SharePointer) []chainoracle.ChainOracleShareRange {
	ranges := make([]chainoracle.ChainOracleShareRange, len(sp.Ranges))
	for i, r := range sp.Ranges {
		ranges[i] = chainoracle.ChainOracleShareRange{
			Start: big.NewInt(int64(r.Start)),
			End:   big.NewInt(int64(r.End)),
		}
	}
	return ranges
}
