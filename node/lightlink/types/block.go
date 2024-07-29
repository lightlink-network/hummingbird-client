package types

import (
	"encoding/binary"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Block struct {
	header       *ethtypes.Header
	uncles       []*ethtypes.Header
	transactions Transactions
	withdrawals  ethtypes.Withdrawals

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header      *ethtypes.Header
	Txs         []*Transaction
	Uncles      []*ethtypes.Header
	Withdrawals []*ethtypes.Withdrawal `rlp:"optional"`
}

// NewBlock creates a new block. The input data is copied, changes to header and to the
// field values will not affect the block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock(header *ethtypes.Header, txs []*Transaction, uncles []*ethtypes.Header, receipts []*ethtypes.Receipt, hasher TrieHasher) *Block {
	b := &Block{header: CopyHeader(header)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = ethtypes.EmptyTxsHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs), hasher)
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = ethtypes.EmptyReceiptsHash
	} else {
		b.header.ReceiptHash = DeriveSha(ethtypes.Receipts(receipts), hasher)
		b.header.Bloom = ethtypes.CreateBloom(receipts)
	}

	if len(uncles) == 0 {
		b.header.UncleHash = ethtypes.EmptyUncleHash
	} else {
		b.header.UncleHash = CalcUncleHash(uncles)
		b.uncles = make([]*ethtypes.Header, len(uncles))
		for i := range uncles {
			b.uncles[i] = CopyHeader(uncles[i])
		}
	}

	return b
}

// NewBlockWithWithdrawals creates a new block with withdrawals. The input data is copied,
// changes to header and to the field values will not affect the block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header are ignored and set to
// values derived from the given txs, uncles and receipts.
func NewBlockWithWithdrawals(header *ethtypes.Header, txs []*Transaction, uncles []*ethtypes.Header, receipts []*ethtypes.Receipt, withdrawals []*ethtypes.Withdrawal, hasher TrieHasher) *Block {
	b := NewBlock(header, txs, uncles, receipts, hasher)

	if withdrawals == nil {
		b.header.WithdrawalsHash = nil
	} else if len(withdrawals) == 0 {
		b.header.WithdrawalsHash = &ethtypes.EmptyWithdrawalsHash
	} else {
		h := DeriveSha(ethtypes.Withdrawals(withdrawals), hasher)
		b.header.WithdrawalsHash = &h
	}

	return b.WithWithdrawals(withdrawals)
}

// CopyHeader creates a deep copy of a block header.
func CopyHeader(h *ethtypes.Header) *ethtypes.Header {
	cpy := *h
	if cpy.Difficulty = new(big.Int); h.Difficulty != nil {
		cpy.Difficulty.Set(h.Difficulty)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if h.BaseFee != nil {
		cpy.BaseFee = new(big.Int).Set(h.BaseFee)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	if h.WithdrawalsHash != nil {
		cpy.WithdrawalsHash = new(common.Hash)
		*cpy.WithdrawalsHash = *h.WithdrawalsHash
	}
	if h.ExcessBlobGas != nil {
		cpy.ExcessBlobGas = new(uint64)
		*cpy.ExcessBlobGas = *h.ExcessBlobGas
	}
	if h.BlobGasUsed != nil {
		cpy.BlobGasUsed = new(uint64)
		*cpy.BlobGasUsed = *h.BlobGasUsed
	}
	if h.ParentBeaconRoot != nil {
		cpy.ParentBeaconRoot = new(common.Hash)
		*cpy.ParentBeaconRoot = *h.ParentBeaconRoot
	}
	return &cpy
}

// DecodeRLP decodes a block from RLP.
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.uncles, b.transactions, b.withdrawals = eb.Header, eb.Uncles, eb.Txs, eb.Withdrawals
	b.size.Store(rlp.ListSize(size))
	return nil
}

// EncodeRLP serializes a block as RLP.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &extblock{
		Header:      b.header,
		Txs:         b.transactions,
		Uncles:      b.uncles,
		Withdrawals: b.withdrawals,
	})
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body struct {
	Transactions []*Transaction
	Uncles       []*ethtypes.Header
	Withdrawals  []*ethtypes.Withdrawal `rlp:"optional"`
}

// Body returns the non-header content of the block.
// Note the returned data is not an independent copy.
func (b *Block) Body() *Body {
	return &Body{b.transactions, b.uncles, b.withdrawals}
}

// Accessors for body data. These do not return a copy because the content
// of the body slices does not affect the cached hash/size in block.

func (b *Block) Uncles() []*ethtypes.Header        { return b.uncles }
func (b *Block) Transactions() Transactions        { return b.transactions }
func (b *Block) Withdrawals() ethtypes.Withdrawals { return b.withdrawals }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

// Header returns the block header (as a copy).
func (b *Block) Header() *ethtypes.Header {
	return CopyHeader(b.header)
}

// Header value accessors. These do copy!

func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64     { return b.header.GasLimit }
func (b *Block) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block) Difficulty() *big.Int { return new(big.Int).Set(b.header.Difficulty) }
func (b *Block) Time() uint64         { return b.header.Time }

func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block) MixDigest() common.Hash   { return b.header.MixDigest }
func (b *Block) Nonce() uint64            { return binary.BigEndian.Uint64(b.header.Nonce[:]) }
func (b *Block) Bloom() ethtypes.Bloom    { return b.header.Bloom }
func (b *Block) Coinbase() common.Address { return b.header.Coinbase }
func (b *Block) Root() common.Hash        { return b.header.Root }
func (b *Block) ParentHash() common.Hash  { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash      { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash { return b.header.ReceiptHash }
func (b *Block) UncleHash() common.Hash   { return b.header.UncleHash }
func (b *Block) Extra() []byte            { return common.CopyBytes(b.header.Extra) }

func (b *Block) BaseFee() *big.Int {
	if b.header.BaseFee == nil {
		return nil
	}
	return new(big.Int).Set(b.header.BaseFee)
}

func (b *Block) BeaconRoot() *common.Hash { return b.header.ParentBeaconRoot }

func (b *Block) ExcessBlobGas() *uint64 {
	var excessBlobGas *uint64
	if b.header.ExcessBlobGas != nil {
		excessBlobGas = new(uint64)
		*excessBlobGas = *b.header.ExcessBlobGas
	}
	return excessBlobGas
}

func (b *Block) BlobGasUsed() *uint64 {
	var blobGasUsed *uint64
	if b.header.BlobGasUsed != nil {
		blobGasUsed = new(uint64)
		*blobGasUsed = *b.header.BlobGasUsed
	}
	return blobGasUsed
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previously cached value.
func (b *Block) Size() uint64 {
	if size := b.size.Load(); size != nil {
		return size.(uint64)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(uint64(c))
	return uint64(c)
}

// SanityCheck can be used to prevent that unbounded fields are
// stuffed with junk data to add processing overhead
func (b *Block) SanityCheck() error {
	return b.header.SanityCheck()
}

// type writeCounter uint64

// func (c *writeCounter) Write(b []byte) (int, error) {
// 	*c += writeCounter(len(b))
// 	return len(b), nil
// }

func CalcUncleHash(uncles []*ethtypes.Header) common.Hash {
	if len(uncles) == 0 {
		return ethtypes.EmptyUncleHash
	}
	return rlpHash(uncles)
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *ethtypes.Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *ethtypes.Header) *Block {
	return &Block{
		header:       CopyHeader(header),
		transactions: b.transactions,
		uncles:       b.uncles,
		withdrawals:  b.withdrawals,
	}
}

// WithBody returns a copy of the block with the given transaction and uncle contents.
func (b *Block) WithBody(transactions []*Transaction, uncles []*ethtypes.Header) *Block {
	block := &Block{
		header:       b.header,
		transactions: make([]*Transaction, len(transactions)),
		uncles:       make([]*ethtypes.Header, len(uncles)),
		withdrawals:  b.withdrawals,
	}
	copy(block.transactions, transactions)
	for i := range uncles {
		block.uncles[i] = CopyHeader(uncles[i])
	}
	return block
}

// WithWithdrawals returns a copy of the block containing the given withdrawals.
func (b *Block) WithWithdrawals(withdrawals []*ethtypes.Withdrawal) *Block {
	block := &Block{
		header:       b.header,
		transactions: b.transactions,
		uncles:       b.uncles,
	}
	if withdrawals != nil {
		block.withdrawals = make([]*ethtypes.Withdrawal, len(withdrawals))
		copy(block.withdrawals, withdrawals)
	}
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks []*Block

// HeaderParentHashFromRLP returns the parentHash of an RLP-encoded
// header. If 'header' is invalid, the zero hash is returned.
func HeaderParentHashFromRLP(header []byte) common.Hash {
	// parentHash is the first list element.
	listContent, _, err := rlp.SplitList(header)
	if err != nil {
		return common.Hash{}
	}
	parentHash, _, err := rlp.SplitString(listContent)
	if err != nil {
		return common.Hash{}
	}
	if len(parentHash) != 32 {
		return common.Hash{}
	}
	return common.BytesToHash(parentHash)
}
