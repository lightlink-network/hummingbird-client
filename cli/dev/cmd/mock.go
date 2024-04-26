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
	"strconv"
	"strings"

	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

type MockData struct {
	RollupHash   common.Hash                                   `json:"rollupHash"`
	RollupHeader canonicalstatechain.CanonicalStateChainHeader `json:"rollupHeader"`
	Headers      []HeaderData                                  `json:"headers"`
	Transactions [][]TransactionData                           `json:"transactions"`
}

type HeaderData struct {
	Header        *utils.L2HeaderJson                 `json:"header"`
	HeaderHash    common.Hash                         `json:"headerHash"`
	ShareProofs   chainoracle.SharesProof             `json:"shareProofs"`
	ShareRanges   []chainoracle.ChainOracleShareRange `json:"shareRanges"`
	PointerProofs []chainoracle.BinaryMerkleProof     `json:"pointerProofs"`
	Shares        [][]byte                            `json:"shares"`
}

type TransactionData struct {
	Transaction   *utils.TxJson                       `json:"transaction"`
	Hash          common.Hash                         `json:"hash"`
	ShareProofs   chainoracle.SharesProof             `json:"shareProofs"`
	ShareRanges   []chainoracle.ChainOracleShareRange `json:"shareRanges"`
	PointerProofs []chainoracle.BinaryMerkleProof     `json:"pointerProofs"`
	Shares        [][]byte                            `json:"shares"`
}

func init() {
	RootCmd.AddCommand(MockCmd)
}

var MockCmd = &cobra.Command{
	Use:     "mock [rblock]:[pointer] [blocks...]",
	Short:   "mock will output mock data for testing using real blocks",
	Long:    "mock will output mock data for testing using real blocks. The first argument is in the form rblock_hash:pointer_index. The Remaining arguments are the indexes of the blocks you want to fetch.",
	Example: "",
	Args:    cobra.MinimumNArgs(3),
	ArgAliases: []string{
		"rblock:pointer",
		"blocks",
	},
	Run: func(cmd *cobra.Command, args []string) {

		// 0. parse args
		rblockPointer := strings.Split(args[0], ":")
		rblockHash := common.HexToHash(rblockPointer[0])
		pointerIndex, _ := strconv.Atoi(rblockPointer[1])

		blocks := []int{}
		for _, block := range args[1:] {
			blockNum, _ := strconv.Atoi(block)
			blocks = append(blocks, blockNum)
		}

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

		// 3. Loop through the blocks and get the header and transaction data
		hds := make([]HeaderData, len(blocks))
		txs := make([][]TransactionData, len(blocks))
		for i, blockNum := range blocks {

			// – Fetch the blocks header data
			hds[i] = getHeaderData(r, rblock, pointerIndex, blockNum)

			// – Fetch the blocks transactions data
			txCount := len(rblock.Bundles[pointerIndex].Blocks[blockNum].Transactions())
			txs[i] = make([]TransactionData, txCount)
			for txNum := 0; txNum < txCount; txNum++ {
				txs[i][txNum] = getTransactionData(r, rblock, pointerIndex, blockNum, txNum)
			}

			log.Info("Got block", "blockNum", blockNum, "headerHash", hds[i].HeaderHash.String(), "txCount", txCount)
		}

		// 4. Output the mock data
		out := &MockData{
			RollupHash:   rblockHash,
			RollupHeader: *rblock.CanonicalStateChainHeader,
			Headers:      hds,
			Transactions: txs,
		}

		printJSON(out)
	},
}

func getHeaderData(r *rollup.Rollup, rblock *rollup.Block, pointerIndex int, blockNum int) HeaderData {
	bundle := rblock.Bundles[pointerIndex]

	// - Get the header
	header := bundle.Blocks[blockNum].Header()
	headerHash := utils.HashHeaderWithoutExtraData(header)

	// - Get the share proofs
	sharePointer, err := bundle.FindHeaderShares(headerHash, r.Namespace())
	panicErr(err, "failed to find header shares")

	shareProof, err := r.Celestia.GetSharesProof(rblock.GetCelestiaPointers()[pointerIndex], sharePointer)
	panicErr(err, "failed to get share proof")

	shareProofs, err := contracts.NewShareProof(shareProof, getAttestations(r.Node, rblock.GetCelestiaPointers()[pointerIndex]))
	panicErr(err, "failed to get share proofs")

	// - Get the block proofs
	blockProofs := node.GetSharesProofs(sharePointer, rblock.Bundles, pointerIndex, r.Namespace())

	return HeaderData{
		Header:        utils.ToL2HeaderJson(header),
		HeaderHash:    headerHash,
		ShareProofs:   *shareProofs,
		Shares:        shares.ToBytes(sharePointer.Shares()),
		ShareRanges:   formatRanges(sharePointer),
		PointerProofs: utils.ToBinaryMerkleProof(blockProofs),
	}
}

func getTransactionData(r *rollup.Rollup, rblock *rollup.Block, pointerIndex int, blockNum int, txNum int) TransactionData {
	bundle := rblock.Bundles[pointerIndex]

	// - Get the Transaction
	tx := bundle.Blocks[blockNum].Transactions()[txNum]

	// - Get the share proofs
	sharePointer, err := bundle.FindTxShares(tx.Hash(), r.Namespace())
	panicErr(err, "failed to find header shares")

	shareProof, err := r.Celestia.GetSharesProof(rblock.GetCelestiaPointers()[pointerIndex], sharePointer)
	panicErr(err, "failed to get share proof")

	shareProofs, err := contracts.NewShareProof(shareProof, getAttestations(r.Node, rblock.GetCelestiaPointers()[pointerIndex]))
	panicErr(err, "failed to get share proofs")

	// - Get the block proofs
	blockProofs := node.GetSharesProofs(sharePointer, rblock.Bundles, pointerIndex, r.Namespace())

	return TransactionData{
		Transaction:   utils.ToTxJson(tx),
		Hash:          tx.Hash(),
		ShareProofs:   *shareProofs,
		Shares:        shares.ToBytes(sharePointer.Shares()),
		ShareRanges:   formatRanges(sharePointer),
		PointerProofs: utils.ToBinaryMerkleProof(blockProofs),
	}
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
