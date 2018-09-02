package main

import (
	bolt "go.etcd.io/bbolt"
	"fmt"
	"github.com/Devs-On-Discord/DoDdy/db"
)

type user struct {
	id string
}

func (u *user) Insert(usersBucket *bolt.Bucket) error {
	userBucket, err := usersBucket.CreateBucket([]byte(u.id))
	if err != nil {
		return fmt.Errorf("create user: %s", err)
	}
	err = userBucket.Put([]byte("id"), []byte(string(u.id)))
	return err
}

func (u *user) Delete() {
	db.Db.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte("users"))
		if usersBucket == nil {
			return nil
		}
		err := usersBucket.DeleteBucket([]byte(u.id))
		if err != nil {
			return fmt.Errorf("delete user: %s", err)
		}
		return nil
	})
}
