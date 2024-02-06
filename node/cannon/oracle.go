package cannon

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

// PreImages is a map of preimage hashes to their values

type PreImages map[common.Hash][]byte

func PreImagesFromProof(proof []string) PreImages {
	newPreImages := make(map[common.Hash][]byte)
	for _, s := range proof {
		ret, _ := hexutil.Decode(s)
		hash := crypto.Keccak256Hash(ret)
		newPreImages[hash] = ret
	}
	return newPreImages
}

func (p PreImages) WriteFile(dir string, hash common.Hash) ([]byte, error) {
	val, ok := p[hash]
	filepath := fmt.Sprintf("%s/%s", dir, hash)
	err := os.WriteFile(filepath, val, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write preimage to file: %v", err)
	}

	expectedHash := crypto.Keccak256Hash(val)
	if ok && expectedHash != hash {
		return nil, fmt.Errorf("corruption for preimage w/ hash: %s", hash)
	}

	return val, nil
}

// PreimageKeyValueWriter wraps the Put method of a backing data store.
// Implements ethdb.KeyValueWriter
type PreimageKeyValueWriter struct {
	preimages PreImages
}

// Put inserts the given value into the key-value data store.
func (w PreimageKeyValueWriter) Put(key []byte, value []byte) error {
	hash := crypto.Keccak256Hash(value)
	if hash != common.BytesToHash(key) {
		panic("bad preimage value write")
	}
	w.preimages[hash] = common.CopyBytes(value)
	return nil
}

// Delete removes the key from the key-value data store.
func (w PreimageKeyValueWriter) Delete(key []byte) error {
	return nil
}

// Oracle will fetch preimages from a remote node

type Oracle struct {
	logger     *slog.Logger
	gethclient *gethclient.Client
	ethclient  *ethclient.Client

	cached map[string]bool
	images PreImages

	inputs  [6]common.Hash
	outputs [2]common.Hash
}

func NewOracle(geth *gethclient.Client, eth *ethclient.Client, logger *slog.Logger) *Oracle {
	return &Oracle{gethclient: geth, ethclient: eth, logger: logger, cached: make(map[string]bool), images: make(PreImages)}
}

func (o *Oracle) Clear() {
	o.cached = nil
	o.images = nil
}

func (o *Oracle) PreFetchStorage(blockNum *big.Int, addr common.Address, skey common.Hash, postProcess func(PreImages)) error {
	// 1. Check if we already have the preimages cached
	cacheKey := fmt.Sprintf("proof_%d_%s_%s", blockNum, addr, skey)
	if o.cached[cacheKey] {
		return nil
	}
	o.cached[cacheKey] = true

	// 2. Otherwise fetch the preimages from the remote node
	newPreImages, err := o.fetchStorage(blockNum, addr, skey)
	if err != nil {
		return err
	}

	// 3. Post process the preimages
	if postProcess != nil {
		postProcess(newPreImages)
	}

	// 4. Cache the preimages
	for k, v := range newPreImages {
		o.images[k] = v
	}

	return nil
}

func (o *Oracle) PreFetchAccount(blockNum *big.Int, addr common.Address, postProcess func(PreImages)) error {
	// 1. Check if we already have the preimages cached
	cacheKey := fmt.Sprintf("proof_%d_%s", blockNum, addr)
	if o.cached[cacheKey] {
		return nil
	}
	o.cached[cacheKey] = true

	// 2. Otherwise fetch the preimages from the remote node
	newPreImages, err := o.fetchAccount(blockNum, addr)
	if err != nil {
		return err
	}

	// 3. Post process the preimages
	if postProcess != nil {
		postProcess(newPreImages)
	}

	// 4. Cache the preimages
	for k, v := range newPreImages {
		o.images[k] = v
	}

	return nil
}

func (o *Oracle) PreFetchCode(blockNum *big.Int, addr common.Address, postProcess func(PreImages)) error {
	// 1. Check if we already have the preimages cached
	cacheKey := fmt.Sprintf("code_%d_%s", blockNum, addr)
	if o.cached[cacheKey] {
		return nil
	}
	o.cached[cacheKey] = true

	// 2. Otherwise fetch the preimages from the remote node
	newPreImages, err := o.fetchCode(blockNum, addr)
	if err != nil {
		return err
	}

	// 3. Post process the preimages
	if postProcess != nil {
		postProcess(newPreImages)
	}

	// 4. Cache the preimages
	for k, v := range newPreImages {
		o.images[k] = v
	}

	return nil
}

func (o *Oracle) PreFetchBlock(blockNum *big.Int, startBlock bool) error {
	images, header, err := o.fetchBlock(blockNum)
	if err != nil {
		return err
	}

	for k, v := range images {
		o.images[k] = v
	}

	// if we are the start block header, just insert the hash as the first input
	if startBlock {
		hash := header.Hash()
		emptyHash := common.Hash{}
		if o.inputs[0] == emptyHash {
			o.inputs[0] = hash
		}
		return nil
	}

	// otherwise if we are the second block
	if header.ParentHash.Cmp(o.inputs[0]) != 0 {
		return fmt.Errorf("block parent incorrect– have: %v, want: %v", header.ParentHash, o.inputs[0])
	}
	o.inputs[1] = header.TxHash
	o.inputs[2] = crypto.Keccak256Hash(header.Coinbase[:])
	o.inputs[3] = header.UncleHash
	o.inputs[4] = common.BigToHash(big.NewInt(int64(header.GasLimit)))
	o.inputs[5] = common.BigToHash(big.NewInt(int64(header.Time)))

	// output – (secret inputs)
	o.outputs[0] = header.Root
	o.outputs[1] = header.ReceiptHash

	// TODO: check txroot is correct
	// TODO: Do we need to save uncle roots? We have no uncles...

	return nil
}

func (o *Oracle) fetchStorage(blockNum *big.Int, addr common.Address, skey common.Hash) (PreImages, error) {
	// 1. get storage proof from remote node
	ctx := context.Background()
	proof, err := o.gethclient.GetProof(ctx, addr, []string{hexutil.Encode(skey[:])}, blockNum)
	if err != nil {
		return nil, err
	}

	// 2. convert storage proof to preimages
	newPreImages := PreImagesFromProof(proof.StorageProof[0].Proof)

	return newPreImages, nil
}

func (o *Oracle) fetchAccount(blockNumber *big.Int, addr common.Address) (PreImages, error) {
	// 1. get account proof from remote node
	ctx := context.Background()
	proof, err := o.gethclient.GetProof(ctx, addr, nil, blockNumber)
	if err != nil {
		return nil, err
	}

	// 2. convert account proof to preimages
	newPreImages := PreImagesFromProof(proof.AccountProof)

	return newPreImages, nil
}

func (o *Oracle) fetchCode(blockNumber *big.Int, addr common.Address) (PreImages, error) {
	// 1. get code from remote node
	ctx := context.Background()
	code, err := o.ethclient.CodeAt(ctx, addr, blockNumber)
	if err != nil {
		return nil, err
	}

	// 2. convert code to preimages
	newPreImages := make(map[common.Hash][]byte)
	hash := crypto.Keccak256Hash(code)
	newPreImages[hash] = code

	return newPreImages, nil
}

func (o *Oracle) fetchBlock(blockNumber *big.Int) (PreImages, *types.Header, error) {
	// 1. get block from remote node
	ctx := context.Background()
	block, err := o.ethclient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// 2. convert block to preimages
	blockHeaderRLP, err := rlp.EncodeToBytes(block.Header())
	if err != nil {
		return nil, nil, err
	}

	newPreImages := make(map[common.Hash][]byte)
	hash := crypto.Keccak256Hash(blockHeaderRLP)
	newPreImages[hash] = blockHeaderRLP

	return newPreImages, nil, nil
}

func (o *Oracle) PreImages() PreImages    { return o.images }
func (o *Oracle) Inputs() [6]common.Hash  { return o.inputs }
func (o *Oracle) Outputs() [2]common.Hash { return o.outputs }

func (o *Oracle) Writer() *PreimageKeyValueWriter {
	return &PreimageKeyValueWriter{preimages: o.images}
}
