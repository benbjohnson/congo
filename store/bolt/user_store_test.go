package bolt_test

import (
	"reflect"
	"testing"

	"github.com/gopheracademy/congo"
	"github.com/gopheracademy/congo/store/bolt"
)

// Ensure the store can create a user.
func TestUserStore_CreateUser(t *testing.T) {
	s := OpenStore()
	defer s.Close()

	u := &congo.User{
		FirstName: "Susy",
		LastName:  "Que",
		Email:     "susy@que.com",
	}

	// Create user & verify ID is returned.
	if err := s.CreateUser(u); err != nil {
		t.Fatal(err)
	} else if u.ID != 1 {
		t.Fatalf("unexpected id: %d", u.ID)
	}

	// Verify user can be retrieved.
	if u, err := s.User(1); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(u, &congo.User{
		ID:        1,
		FirstName: "Susy",
		LastName:  "Que",
		Email:     "susy@que.com",
	}) {
		t.Fatalf("unexpected user: %#v", u)
	}
}

// Ensure the store returns an error when creating a user without an email.
func TestUserStore_CreateUser_ErrEmailRequired(t *testing.T) {
	s := OpenStore()
	defer s.Close()
	if err := s.CreateUser(&congo.User{Email: ""}); err != bolt.ErrEmailRequired {
		t.Fatalf("unexpected error: %s", err)
	}
}

// Ensure the store returns a nil user if no email matches.
func TestUserStore_UserByEmail_NotFound(t *testing.T) {
	s := OpenStore()
	defer s.Close()

	// Verify no error or user is returned.
	if u, err := s.UserByEmail("no_such_user"); err != nil {
		t.Fatal(err)
	} else if u != nil {
		t.Fatalf("unexpected user: %#v", u)
	}
}

// Ensure the store can retrieve a user by email.
func TestUserStore_UserByEmail(t *testing.T) {
	s := OpenStore()
	defer s.Close()

	// Create user.
	if err := s.CreateUser(&congo.User{Email: "susy@que.com"}); err != nil {
		t.Fatal(err)
	}

	// Verify user can be retrieved by email.
	if u, err := s.UserByEmail("susy@que.com"); err != nil {
		t.Fatal(err)
	} else if u == nil || u.ID != 1 {
		t.Fatalf("unexpected user: %#v", u)
	}
}

// Ensure the store can delete a user.
func TestUserStore_DeleteUser(t *testing.T) {
	s := OpenStore()
	defer s.Close()

	// Create user.
	if err := s.CreateUser(&congo.User{Email: "susy@que.com"}); err != nil {
		t.Fatal(err)
	}

	// Remove user.
	if err := s.DeleteUser(1); err != nil {
		t.Fatal(err)
	}

	// Verify user cannot be retrieved.
	if u, err := s.User(1); err != nil {
		t.Fatal(err)
	} else if u != nil {
		t.Fatalf("unexpected user: %#v", u)
	}
}

// Ensure the store returns an error if deleting a non-existent user.
func TestUserStore_DeleteUser_ErrUserNotFound(t *testing.T) {
	s := OpenStore()
	defer s.Close()
	if err := s.DeleteUser(100); err != bolt.ErrUserNotFound {
		t.Fatalf("unexpected error: %s", err)
	}
}
