package cmd

import (
	"fmt"
	"math/big"

	"github.com/celestiaorg/go-square/shares"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"

	"hummingbird/node"
	"hummingbird/node/contracts"
	chainoracleContract "hummingbird/node/contracts/ChainOracle.sol"
	"hummingbird/rollup"
	"hummingbird/utils"
)

func init() {
	FetchCmd.Flags().StringVar(&format, "format", "json", "format of the output (json, pretty)")
	FetchCmd.Flags().BoolVar(&withProof, "proof", false, "whether to generate share proofs in the output")
	FetchCmd.Flags().BoolVar(&checkProof, "check-proof", false, "whether to check the proof")
	FetchCmd.Flags().BoolVar(&checkPointer, "check-pointer", false, "whether to check the pointer")

	// Add the fetch command to the root command
	RootCmd.AddCommand(FetchCmd)
}

var (
	// variables for the fetch command
	format       string // format of the output (json)
	withProof    bool   // whether to generate share proofs in the output
	checkProof   bool   // whether to check the proof
	checkPointer bool   // whether to check the pointer

	// fetch command
	FetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetch will fetch an item (either: header or tx) from a given rollup block",
		Long:  "fetch will fetch an item (either: header or tx) from a given rollup block. This can be used for generating test data for the smart contracts.",
		Args:  cobra.MinimumNArgs(3),
		ArgAliases: []string{
			"data-type",
			"rblock",
			"data-hash",
		},
		Run: func(cmd *cobra.Command, args []string) {

			// 0. parse args
			dataType := args[0]
			rblockHash := common.HexToHash(args[1])
			dataHash := common.HexToHash(args[2])

			// 1. make node
			n, log, err := makeNode()
			panicErr(err, "failed to create node")
			r := rollup.NewRollup(n, &rollup.Opts{
				Logger: log.With("ctx", "Rollup"),
			})

			// 2. Get rblock and celestia pointer
			rblock, err := r.GetBlockByHash(rblockHash)
			panicErr(err, "failed to get rollup block")
			log.Debug("✔️  Got Rollup Block", "hash", rblockHash)

			// 3. Fetch the item
			var item any
			var sharePointer *node.SharePointer
			var celPointer *node.CelestiaPointer
			var pointerIndex uint8

			switch dataType {
			case "header":
				sharePointer, pointerIndex, err = node.FindHeaderSharesInBundles(rblock.Bundles, dataHash, n.Namespace())
				celPointer = rblock.GetCelestiaPointers()[int(pointerIndex)]
				panicErr(err, "failed to find header shares")
				h := &types.Header{}
				err = rlp.DecodeBytes(sharePointer.Bytes(), h)
				panicErr(err, "failed to decode header")
				item = &h
			case "tx":
				sharePointer, pointerIndex, err = node.FindTxSharesInBundles(rblock.Bundles, dataHash, n.Namespace())
				celPointer = rblock.GetCelestiaPointers()[int(pointerIndex)]
				panicErr(err, "failed to find tx shares")
				tx := &types.Transaction{}
				err = rlp.DecodeBytes(sharePointer.Bytes(), tx)
				panicErr(err, "failed to decode tx")
				item = &tx
			default:
				panicErr(err, "invalid data type")
			}
			log.Debug("✔️  Fetched Item", "type", dataType)

			if checkPointer {
				log.Warn("Not implemented: check-pointer")
				// ss, err := n.Celestia.GetShares(&node.CelestiaPointer{
				// 	Height:     celPointer.Height,
				// 	ShareStart: uint64(pointer.StartShare),
				// 	ShareLen:   uint64(pointer.EndShare()-pointer.StartShare) + 1,
				// })
				// panicErr(err, "failed to get shares")

				// for i, s := range pointer.Shares() {
				// 	if !bytes.Equal(s.ToBytes(), ss[i].ToBytes()) {
				// 		log.Info("share pointer length", "expect", uint64(pointer.EndShare()-pointer.StartShare), "got", len(pointer.Ranges))
				// 		log.Error("share does not match", "index", i, "expected", fmt.Sprintf("%x", s), "got", fmt.Sprintf("%x", ss[i]))
				// 		panic("check-pointer: share does not match")
				// 	}
				// 	log.Debug("✔️  Share matches", "index", i)
				// }
			}

			// 5. Generate proofs
			var proof *chainoracleContract.SharesProof
			if withProof {
				shareProof, err := n.Celestia.GetSharesProof(celPointer, sharePointer)
				panicErr(err, "failed to get share proof")

				commitment, err := n.Ethereum.GetBlobstreamCommitment(int64(celPointer.Height))
				panicErr(err, "failed to get blobstream commitment")

				celProof, err := n.Celestia.GetProof(celPointer, commitment.StartBlock, commitment.EndBlock, *commitment.ProofNonce)
				panicErr(err, "failed to get celestia proof")

				attestationProof := chainoracleContract.AttestationProof{
					TupleRootNonce: celProof.Nonce,
					Tuple: chainoracleContract.DataRootTuple{
						Height:   celProof.Tuple.Height,
						DataRoot: celProof.Tuple.DataRoot,
					},
					Proof: chainoracleContract.BinaryMerkleProof{
						SideNodes: celProof.WrappedProof.SideNodes,
						Key:       celProof.WrappedProof.Key,
						NumLeaves: celProof.WrappedProof.NumLeaves,
					},
				}

				proof, err = contracts.NewShareProof(shareProof, attestationProof)
				panicErr(err, "failed to create proof")

				// check the proof can be decoded
				if checkProof {
					ss, err := utils.BytesToShares(proof.Data)
					panicErr(err, "failed to convert proof to shares")
					switch dataType {
					case "header":
						decH, err := sharesToHeader(ss, sharePointer.Ranges)
						panicErr(err, "failed to convert shares to header")

						decH.Extra = common.Hex2Bytes("0x")
						if decH.Hash().Hex() != dataHash.Hex() {
							panicErr(err, "proof does not match data")
						}

					default:
						log.Warn("proof check not implemented for tx")
					}

					fmt.Println("✔️  Proof is valid")
				}
			}

			// 5. Generate output
			ranges := []chainoracleContract.ChainOracleShareRange{}
			for _, r := range sharePointer.Ranges {
				ranges = append(ranges, chainoracleContract.ChainOracleShareRange{
					Start: big.NewInt(int64(r.Start)),
					End:   big.NewInt(int64(r.End)),
				})
			}

			shareBytes := [][]byte{}
			for _, s := range sharePointer.Shares() {
				shareBytes = append(shareBytes, s.ToBytes())
			}

			output := &Output[any]{
				RBlock: rblockHash,
				Hash:   dataHash,
				Data:   item,
				Shares: shareBytes,
				Proof:  proof,
				Ranges: ranges,
			}

			// 6. Print output
			switch format {
			case "json":
				printJSON(output)
			default:
				printPretty(output)
			}
		},
	}
)

type Output[T any] struct {
	RBlock common.Hash                                 `json:"rblock"`
	Hash   common.Hash                                 `json:"hash"`
	Data   T                                           `json:"content"`
	Shares [][]byte                                    `json:"shares,omitempty"`
	Ranges []chainoracleContract.ChainOracleShareRange `json:"ranges,omitempty"`
	Proof  *chainoracleContract.SharesProof            `json:"proof,omitempty"`
}

func sharesToHeader(s []shares.Share, ranges []node.ShareRange) (*types.Header, error) {
	data := []byte{}
	for i, r := range ranges {
		data = append(data, s[i].ToBytes()[r.Start:r.End]...)
	}

	header := &types.Header{}
	return header, rlp.DecodeBytes(data, &header)
}
