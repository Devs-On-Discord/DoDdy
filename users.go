package main

import (
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
)

//TODO: remove users on disconnect
type users struct {
	user map[uint64]user
}

func (u *users) load(id uint64) user {
	existingUser, exists := u.user[id]
	if exists {
		return existingUser
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user-" + string(id)))
		if b == nil {
			existingUser := user{id: id}
			err := existingUser.Insert(tx)
			u.user[id] = existingUser
			return err
		}
		id := b.Get([]byte("id"))
		existingUser = user{id: binary.BigEndian.Uint64(id)}
		return nil
	})
	return existingUser

}
