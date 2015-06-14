package bolt_test

import (
	"io/ioutil"
	"os"

	"github.com/gopheracademy/congo/store/bolt"
)

// Store represents a test wrapper for bolt.Store.
type Store struct {
	*bolt.Store
}

// NewStore returns a new instance of Store in a temporary path.
func NewStore() *Store {
	path, _ := ioutil.TempDir("", "congo-bolt-")
	return &Store{bolt.NewStore(path)}
}

// OpenStore returns an open instance of Store.
func OpenStore() *Store {
	s := NewStore()
	if err := s.Open(); err != nil {
		panic(err)
	}
	return s
}

// Close closes and deletes the database.
func (s *Store) Close() {
	s.Store.Close()
	os.RemoveAll(s.Path())
}
