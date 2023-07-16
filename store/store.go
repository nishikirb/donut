package store

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
)

// // Store is an interface for storing and retrieving data.
// type Store interface {
// 	Get(bucket string, key string, value any) error
// 	Set(bucket string, key string, value any) error
// }

// BoltStore is a Store implementation that uses BoltDB.
type BoltStore struct {
	db *bolt.DB
}

var (
	store = &BoltStore{}
	once  sync.Once
)

// Open opens a BoltDB database.
func Open(file string, buckets []string) (*BoltStore, error) {
	if err := os.MkdirAll(file, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore{
		db: db,
	}, nil
}

// Init initializes the store.
func Init(file string, buckets []string) error {
	var err error
	once.Do(func() {
		store, err = Open(file, buckets)
	})
	return err
}

// Get retrieves a value from the store.
func Get(bucket string, key string, value any) error {
	return store.Get(bucket, key, value)
}

// Get retrieves a value from the store.
func (s *BoltStore) Get(bucket string, key string, value any) error {
	var raw []byte
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		raw = slices.Clone(b.Get([]byte(key)))
		return nil
	}); err != nil {
		return err
	} else if raw == nil {
		return nil
	}
	return json.Unmarshal(raw, value)
}

// Set stores a value in the store.
func Set(bucket string, key string, value any) error {
	return store.Set(bucket, key, value)
}

// Set stores a value in the store.
func (s *BoltStore) Set(bucket string, key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), raw)
	}); err != nil {
		return err
	}
	return nil
}

// Close closes the store.
func Close() error {
	return store.db.Close()
}
