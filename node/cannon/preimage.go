package cannon

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
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

type PreImageOracle struct {
	logger     *slog.Logger
	gethclient *gethclient.Client
	ethclient  *ethclient.Client

	cached map[string]bool
	images PreImages
}

func NewPreImageOracle(client *gethclient.Client) *PreImageOracle {
	return &PreImageOracle{gethclient: client}
}

func (o *PreImageOracle) Clear() {
	o.cached = nil
	o.images = nil
}

func (o *PreImageOracle) PreFetchStorage(blockNum *big.Int, addr common.Address, skey common.Hash, postProcess func(PreImages)) error {
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

func (o *PreImageOracle) PreFetchAccount(blockNum *big.Int, addr common.Address, postProcess func(PreImages)) error {
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

func (o *PreImageOracle) PreFetchCode(blockNum *big.Int, addr common.Address, postProcess func(PreImages)) error {
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

func (o *PreImageOracle) fetchStorage(blockNum *big.Int, addr common.Address, skey common.Hash) (PreImages, error) {
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

func (o *PreImageOracle) fetchAccount(blockNumber *big.Int, addr common.Address) (PreImages, error) {
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

func (o *PreImageOracle) fetchCode(blockNumber *big.Int, addr common.Address) (PreImages, error) {
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

func (o *PreImageOracle) PreImages() PreImages {
	return o.images
}

func (o *PreImageOracle) Writer() *PreimageKeyValueWriter {
	return &PreimageKeyValueWriter{preimages: o.images}
}
