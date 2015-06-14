package bolt

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

//go:generate protoc --gogo_out=. internal/internal.proto

// Store represents the BoltDB-backed data store.
type Store struct {
	path string
	db   *bolt.DB

	*UserStore
}

// NewStore returns a new instance of Store at the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Path returns the path the database was initialized with.
func (s *Store) Path() string { return s.path }

// Open opens and initializes the database.
func (s *Store) Open() error {
	// Create path if not exists.
	if err := os.MkdirAll(s.path, 0777); err != nil {
		return err
	}

	// Open underlying bolt database.
	db, err := bolt.Open(filepath.Join(s.path, "db"), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.db = db

	// Create individual stores.
	s.UserStore = &UserStore{db: db}

	// Initialize stores.
	if err := func() error {
		if err := s.UserStore.init(); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		s.Close()
		return err
	}

	return nil
}

// Close shuts down the database.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func u64tob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btou64(b []byte) uint64 { return binary.BigEndian.Uint64(b) }
