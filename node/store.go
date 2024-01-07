package node

import "github.com/syndtr/goleveldb/leveldb"

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
}

type LDBStore struct {
	db *leveldb.DB
}

func NewLDBStore(path string) (*LDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
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
