package node

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	GetDAPointer(hash common.Hash) (*CelestiaPointer, error)
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
