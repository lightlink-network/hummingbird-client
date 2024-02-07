package cannon

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/pellartech/minigeth/common"
	"github.com/pellartech/minigeth/consensus/misc"
	"github.com/pellartech/minigeth/core"
	"github.com/pellartech/minigeth/core/state"
	"github.com/pellartech/minigeth/core/types"
	"github.com/pellartech/minigeth/core/vm"
	"github.com/pellartech/minigeth/crypto"
	"github.com/pellartech/minigeth/oracle"
	"github.com/pellartech/minigeth/params"
	"github.com/pellartech/minigeth/rlp"
	"github.com/pellartech/minigeth/trie"
)

type MiniGethOpts struct {
	Logger  *slog.Logger
	NodeURL string // URL of the EVM compatible node
	BaseDir string // Directory to store the witness data, defaults to /tmp/cannon
}

// MiniGeth is a wrapper for github.com/pellartech/minigeth. It is used to
// fetch preimages for a given block.
type MiniGeth struct {
	Opts *MiniGethOpts
}

func NewMiniGeth(opts *MiniGethOpts) *MiniGeth {
	return &MiniGeth{Opts: opts}
}

func (m *MiniGeth) Process(blockNum int64) error {
	m.Opts.Logger.Info("Processing block", "blockNum", blockNum)

	pkw := oracle.PreimageKeyValueWriter{}
	pkwtrie := trie.NewStackTrie(pkw)

	// 1. Set up the oracle
	oracle.SetNodeUrl(m.Opts.NodeURL)
	oracle.SetRoot(fmt.Sprintf("%s/0_%d", m.Opts.BaseDir, blockNum))

	// 2. Prefetch the block and the next block
	m.Opts.Logger.Info("Prefetching blocks", "start", blockNum, "end", blockNum+1)
	oracle.PrefetchBlock(big.NewInt(blockNum), true, nil)
	oracle.PrefetchBlock(big.NewInt(blockNum+1), false, pkwtrie)
	// 3. Commit those transactions
	hash, err := pkwtrie.Commit()
	if err != nil {
		return err
	}
	m.Opts.Logger.Info("Committed Transactions", "hash", hash)

	// init secp256k1BytePoints?
	crypto.S256()

	// 4. Get inputs
	inputBytes := oracle.Preimage(oracle.InputHash())
	var inputs [6]common.Hash
	for i := 0; i < len(inputs); i++ {
		inputs[i] = common.BytesToHash(inputBytes[i*0x20 : i*0x20+0x20])
	}

	// 3. read start block header
	var parent types.Header
	err = rlp.DecodeBytes(oracle.Preimage(inputs[0]), &parent)
	if err != nil {
		return err
	}

	// 4. read new header
	var newheader types.Header
	// - from parent
	newheader.ParentHash = parent.Hash()
	newheader.Number = big.NewInt(0).Add(parent.Number, big.NewInt(1))
	newheader.BaseFee = misc.CalcBaseFee(params.MainnetChainConfig, &parent)
	// - from input oracle
	newheader.TxHash = inputs[1]
	newheader.Coinbase = common.BigToAddress(inputs[2].Big())
	newheader.UncleHash = inputs[3]
	newheader.GasLimit = inputs[4].Big().Uint64()
	newheader.Time = inputs[5].Big().Uint64()

	// 5. Process the state transition
	bc := core.NewBlockChain(&parent)
	database := state.NewDatabase(parent)
	statedb, _ := state.New(parent.Root, database, nil)
	vmconfig := vm.Config{}
	processor := core.NewStateProcessor(params.MainnetChainConfig, bc, bc.Engine())
	m.Opts.Logger.Info("Processing state", "from", parent.Number, "to", newheader.Number)

	newheader.Difficulty = bc.Engine().CalcDifficulty(bc, newheader.Time, &parent)

	// - read transactions
	var txs []*types.Transaction
	triedb := trie.NewDatabase(parent)
	tt, _ := trie.New(newheader.TxHash, &triedb)
	tni := tt.NodeIterator([]byte{})
	for tni.Next(true) {
		if tni.Leaf() {
			tx := types.Transaction{}
			var rlpKey uint64
			err = rlp.DecodeBytes(tni.LeafKey(), &rlpKey)
			if err != nil {
				return err
			}
			err = tx.UnmarshalBinary(tni.LeafBlob())
			if err != nil {
				return err
			}
			// TODO: resize an array in go?
			for uint64(len(txs)) <= rlpKey {
				txs = append(txs, nil)
			}
			txs[rlpKey] = &tx
		}
	}
	m.Opts.Logger.Info("Read transactions", "count", len(txs))
	// TODO: OMG the transaction ordering isn't fixed

	// - read uncles
	var uncles []*types.Header
	err = rlp.DecodeBytes(oracle.Preimage(newheader.UncleHash), &uncles)
	if err != nil {
		return err
	}

	// - create block
	var receipts []*types.Receipt
	block := types.NewBlock(&newheader, txs, uncles, receipts, trie.NewStackTrie(nil))
	m.Opts.Logger.Info("made block", "parent", newheader.ParentHash)

	// - check block
	if newheader.TxHash != block.Header().TxHash {
		return fmt.Errorf("wrong txs for block")
	}
	if newheader.UncleHash != block.Header().UncleHash {
		return fmt.Errorf("wrong uncles for block %s %s", newheader.UncleHash, block.Header().UncleHash)
	}

	// validateState is more complete, gas used + bloom also
	receipts, _, _, err = processor.Process(block, statedb, vmconfig)
	receiptSha := types.DeriveSha(types.Receipts(receipts), trie.NewStackTrie(nil))
	if err != nil {
		return fmt.Errorf("error processing block: %s", err)
	}
	newRoot := statedb.IntermediateRoot(bc.Config().IsEIP158(newheader.Number))

	// 6. Write the new state
	m.Opts.Logger.Debug("Writing new state", "root", newRoot)
	m.Opts.Logger.Debug("receipts", "count", len(receipts), "hash", receiptSha)
	m.Opts.Logger.Info("process done", "from", parent.Root, "to", newRoot)
	oracle.Output(newRoot, receiptSha)

	return nil
}
