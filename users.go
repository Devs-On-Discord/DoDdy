package main

import (
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
	"fmt"
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
		usersBucket, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("could not create or get users bucket: %s", err)
		}
		userBucket := usersBucket.Bucket([]byte("user-" + string(id)))
		if userBucket == nil {
			existingUser := user{id: id}
			err := existingUser.Insert(usersBucket)
			u.user[id] = existingUser
			return err
		}
		id := usersBucket.Get([]byte("id"))
		existingUser = user{id: binary.BigEndian.Uint64(id)}
		return nil
	})
	return existingUser
}
