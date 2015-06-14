package bolt

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/gopheracademy/congo"
	"github.com/gopheracademy/congo/store/bolt/internal"
)

var (
	// ErrUserNotFound is returned when mutating a user that doesn't exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailRequired is returned when creating a user without an email.
	ErrEmailRequired = errors.New("email address required")
)

// UserStore represents the data store for user data.
type UserStore struct {
	db *bolt.DB
}

// initialize creates the top-level buckets for the store.
func (s *UserStore) init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		tx.CreateBucketIfNotExists([]byte("users.email"))
		return nil
	})
}

// User returns a user by ID.
func (s *UserStore) User(id uint64) (u *congo.User, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("users")).Get(u64tob(id)); v != nil {
			u, err = internal.UnmarshalUser(v)
			return err
		}
		return nil
	})
	return
}

// UserByEmail returns a user by email address.
func (s *UserStore) UserByEmail(email string) (u *congo.User, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		// Lookup user ID from index.
		key := tx.Bucket([]byte("users.email")).Get([]byte(email))
		if key == nil {
			return nil
		}

		// Retrieve the user by ID.
		if v := tx.Bucket([]byte("users")).Get(key); v != nil {
			u, err = internal.UnmarshalUser(v)
			return err
		}
		return nil
	})
	return
}

// CreateUser creates a user and returns a new instance with its ID.
// The new identifier is set to u.ID.
func (s *UserStore) CreateUser(u *congo.User) error {
	// Ensure user has an email address.
	if u.Email == "" {
		return ErrEmailRequired
	}

	// Insert into database.
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		// Generate autoincrementing ID.
		u.ID, _ = b.NextSequence()

		// Encode user.
		buf, err := internal.MarshalUser(u)
		if err != nil {
			return err
		}

		// Insert encoded user by ID.
		if err := b.Put(u64tob(u.ID), buf); err != nil {
			return err
		}

		// Insert email into index.
		if err := tx.Bucket([]byte("users.email")).Put([]byte(u.Email), u64tob(u.ID)); err != nil {
			return err
		}

		return nil
	})
}

// DeleteUser removes a user by ID.
func (s *UserStore) DeleteUser(id uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// Retrieve user.
		v := tx.Bucket([]byte("users")).Get(u64tob(id))
		if v == nil {
			return ErrUserNotFound
		}

		// Decode user.
		u, err := internal.UnmarshalUser(v)
		if err != nil {
			return err
		}

		// Remove from bucket.
		if err := tx.Bucket([]byte("users")).Delete(u64tob(id)); err != nil {
			return err
		}
		if err := tx.Bucket([]byte("users.email")).Delete([]byte(u.Email)); err != nil {
			return err
		}

		return nil
	})
}
