package store

import "github.com/gopheracademy/congo"

// Interface is a super interface of all substores.
type Interface interface {
	UserStore
}

// UserStore represents the interface for user storage.
type UserStore interface {
	User(id int) (*congo.User, error)
	UserByEmail(email string) (*congo.User, error)
	CreateUser(u *congo.User) error
	DeleteUser(id uint64) error
}
