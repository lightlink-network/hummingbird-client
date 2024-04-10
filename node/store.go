package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	GetDAPointer(hash common.Hash) (*CelestiaPointer, error)
	PutBundle(bundle *Bundle) error
	GetBundle(startBlock uint64, endBlock uint64) (*Bundle, error)
}

type LDBStore struct {
	db *leveldb.DB
}

func NewLDBStore(path string) (*LDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %w", err)
	}

	return &LDBStore{db: db}, nil
}

func (l *LDBStore) Get(key []byte) ([]byte, error) {
	return l.db.Get(key, nil)
}

func (l *LDBStore) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *LDBStore) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

func (l *LDBStore) GetDAPointer(hash common.Hash) (*CelestiaPointer, error) {
	if l.db == nil {
		return nil, errors.New("no store")
	}

	key := append([]byte("pointer_"), hash[:]...)
	buf, err := l.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get celestia pointer from store: %w", err)
	}

	pointer := &CelestiaPointer{}
	err = json.Unmarshal(buf, pointer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal celestia pointer: %w", err)
	}

	return pointer, nil
}

func (l *LDBStore) PutBundle(bundle *Bundle) error {
	if l.db == nil {
		return errors.New("no store")
	}

	startL2height := strconv.FormatUint(bundle.Blocks[0].NumberU64(), 10)
	endL2height := strconv.FormatUint(bundle.Height(), 10)
	key := []byte("bundle_" + startL2height + endL2height)

	buf, err := bundle.EncodeRLP()
	if err != nil {
		return fmt.Errorf("failed to marshal bundle: %w", err)
	}

	if err := l.Put(key, buf); err != nil {
		return fmt.Errorf("createNextBlock: Failed to store header: %w", err)
	}

	return l.Put(key, buf)
}

func (l *LDBStore) GetBundle(startBlock uint64, endBlock uint64) (*Bundle, error) {
	if l.db == nil {
		return nil, errors.New("no store")
	}

	startL2height := strconv.FormatUint(startBlock, 10)
	endL2height := strconv.FormatUint(endBlock, 10)
	key := []byte("bundle_" + startL2height + endL2height)

	buf, err := l.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle from store: %w", err)
	}

	bundle := &Bundle{}
	err = bundle.DecodeRLP(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle: %w", err)
	}

	return bundle, nil
}
