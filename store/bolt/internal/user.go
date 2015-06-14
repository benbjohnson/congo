package internal

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gopheracademy/congo"
)

// MarshalUser encodes u into a binary form.
func MarshalUser(u *congo.User) ([]byte, error) {
	return proto.Marshal(encodeUser(u))
}

// UnmarshalUser decodes data into a user.
func UnmarshalUser(data []byte) (*congo.User, error) {
	var u User
	if err := proto.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return decodeUser(&u), nil
}

// encodeUser encodes a user into its internal representation.
func encodeUser(u *congo.User) *User {
	return &User{
		ID:        proto.Uint64(u.ID),
		FirstName: proto.String(u.FirstName),
		LastName:  proto.String(u.LastName),
		Email:     proto.String(u.Email),
	}
}

// decodeUser decodes a user from its internal representation.
func decodeUser(pb *User) *congo.User {
	return &congo.User{
		ID:        pb.GetID(),
		FirstName: pb.GetFirstName(),
		LastName:  pb.GetLastName(),
		Email:     pb.GetEmail(),
	}
}
