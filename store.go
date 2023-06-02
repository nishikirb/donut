package donut

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
)

type Store struct {
	db *bolt.DB
}

var (
	store     *Store
	initStore sync.Once
)

const FileEntryBucket = "file_entries"

func defaultDBFile() string {
	return filepath.Join(defaultStateDir(), "donut.db")
}

func InitStore() error {
	var e error
	initStore.Do(func() {
		if err := os.MkdirAll(filepath.Dir(defaultDBFile()), os.ModePerm); err != nil {
			e = err
		}

		var err error
		db, err := bolt.Open(defaultDBFile(), 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			e = err
		}

		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(FileEntryBucket))
			return err
		}); err != nil {
			e = err
		}

		store = &Store{
			db: db,
		}
	})
	return e
}

func GetStore() *Store {
	return store
}

func (s *Store) Get(bucket string, key string, value any) error {
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

func (s *Store) Set(bucket string, key string, value any) error {
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

func (s *Store) Close() error {
	return s.db.Close()
}
