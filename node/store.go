package node

import (
	"encoding/json"
	"errors"
	"fmt"

	challengeContract "hummingbird/node/contracts/Challenge.sol"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	GetDAPointer(hash common.Hash) (*CelestiaPointer, error)
	StoreActiveDAChallenge(c *challengeContract.ChallengeChallengeDAUpdate) error
	GetActiveDAChallenges() ([]*challengeContract.ChallengeChallengeDAUpdate, error)
	DeleteActiveDAChallenge(blockHash common.Hash) error
	StoreLastScannedBlockNumber(blockNumber uint64) error
	GetLastScannedBlockNumber() (uint64, error)
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

func (l *LDBStore) StoreActiveDAChallenge(c *challengeContract.ChallengeChallengeDAUpdate) error {
	if l.db == nil {
		return errors.New("no store")
	}

	key := append([]byte("da_challenge_"), c.BlockHash[:]...)
	buf, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal challenge: %w", err)
	}

	err = l.Put(key, buf)
	if err != nil {
		return fmt.Errorf("failed to store challenge: %w", err)
	}

	return nil
}

func (l *LDBStore) GetActiveDAChallenges() ([]*challengeContract.ChallengeChallengeDAUpdate, error) {
	if l.db == nil {
		return nil, errors.New("no store")
	}

	iter := l.db.NewIterator(nil, nil)
	defer iter.Release()

	challenges := []*challengeContract.ChallengeChallengeDAUpdate{}
	for iter.Next() {
		if string(iter.Key()[:12]) == "da_challenge" {
			challenge := &challengeContract.ChallengeChallengeDAUpdate{}
			err := json.Unmarshal(iter.Value(), challenge)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal challenge: %w", err)
			}

			if challenge.Status == 1 {
				challenges = append(challenges, challenge)
			}
		}
	}

	return challenges, nil
}

func (l *LDBStore) DeleteActiveDAChallenge(blockHash common.Hash) error {
	if l.db == nil {
		return errors.New("no store")
	}

	key := append([]byte("da_challenge_"), blockHash[:]...)
	err := l.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete challenge: %w", err)
	}

	return nil
}

func (l *LDBStore) StoreLastScannedBlockNumber(blockNumber uint64) error {
	if l.db == nil {
		return errors.New("no store")
	}

	buf, err := json.Marshal(blockNumber)
	if err != nil {
		return fmt.Errorf("failed to marshal block number: %w", err)
	}

	err = l.Put([]byte("last_scanned_block_number"), buf)
	if err != nil {
		return fmt.Errorf("failed to store block number: %w", err)
	}

	return nil
}

func (l *LDBStore) GetLastScannedBlockNumber() (uint64, error) {
	if l.db == nil {
		return 0, errors.New("no store")
	}

	buf, err := l.Get([]byte("last_scanned_block_number"))
	if err != nil {
		return 0, fmt.Errorf("failed to get block number from store: %w", err)
	}

	var blockNumber uint64
	err = json.Unmarshal(buf, &blockNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	return blockNumber, nil
}
