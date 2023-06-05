package donut

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
)

type Store struct {
	db *bolt.DB
}

const FileEntryBucket = "file_entries"

func OpenStore() (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(defaultDBFile()), os.ModePerm); err != nil {
		return nil, err
	}

	db, err := bolt.Open(defaultDBFile(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(FileEntryBucket))
		return err
	}); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
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

func defaultDBFile() string {
	return filepath.Join(defaultStateDir(), "donut.db")
}
