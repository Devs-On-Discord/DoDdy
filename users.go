package main

import (
	bolt "go.etcd.io/bbolt"
	"encoding/binary"
)

type users struct {
	user map[uint64]user
}

func (u *users) load(id uint64) (user) {
	if existingUser, exists := u.user[id]; exists {
		return existingUser
	} else {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("user-" + string(id)))
			if b == nil {
				existingUser := user{id: id}
				existingUser.Insert()
				u.user[id] = existingUser
				return nil
			} else {
				id := b.Get([]byte("id"))
				existingUser = user{id: binary.BigEndian.Uint64(id)}
			}
			return nil
		})
		return existingUser
	}
}
