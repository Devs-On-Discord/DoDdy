package main

import (
	"fmt"

	"github.com/Devs-On-Discord/DoDdy/db"
	bolt "go.etcd.io/bbolt"
)

//TODO: remove users on disconnect
type users struct {
	user map[string]user
}

func (u *users) load(id string) user {
	existingUser, exists := u.user[id]
	if exists {
		return existingUser
	}
	db.DB.Update(func(tx *bolt.Tx) error {
		usersBucket, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("could not create or get users bucket: %s", err)
		}
		userBucket := usersBucket.Bucket([]byte(id))
		if userBucket == nil {
			existingUser := user{id: id}
			err := existingUser.Insert(usersBucket)
			u.user[id] = existingUser
			return err
		}
		id := usersBucket.Get([]byte("id"))
		existingUser = user{id: string(id)}
		return nil
	})
	return existingUser
}
