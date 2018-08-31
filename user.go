package main

import "fmt"
import bolt "go.etcd.io/bbolt"

type user struct {
	id uint64
}

func (u *user) Insert(tx *bolt.Tx) (error) {
	b, err := tx.CreateBucket([]byte("user-" + string(u.id)))
	if err != nil {
		return fmt.Errorf("create user: %s", err)
	}
	err = b.Put([]byte("id"), []byte(string(u.id)))
	return nil
}

func (u *user) Delete() {
	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("user-" + string(u.id)))
		if err != nil {
			return fmt.Errorf("delete user: %s", err)
		}
		return nil
	})
}
