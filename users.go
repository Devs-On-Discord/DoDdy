package main

/*import (
	"fmt"

	"github.com/Devs-On-Discord/DoDdy/db"
	bolt "go.etcd.io/bbolt"
)

//TODO: remove users on disconnect
type users struct {
	db    *bolt.DB
	users map[string]*user
}

type user struct {
	id string
}

func (u *users) Init(db *bolt.DB) {
	u.db = db
	u.users = make(map[string]*user)
}

func (u *users) load(id string) *user {
	if existingUser, exists := u.users[id]; exists {
		return existingUser
	} else {
		db.DB.View(func(tx *bolt.Tx) error {
			usersBucket, err := tx.CreateBucketIfNotExists([]byte("users"))
			if err != nil {
				return fmt.Errorf("could not create or get users bucket: %s", err)
			}
			userBucket := usersBucket.Bucket([]byte(id))
			if userBucket == nil {
				existingUser := &user{id: id}
				err := existingUser.Insert(usersBucket)
				u.users[id] = existingUser
				return err
			}
			id := usersBucket.Get([]byte("id"))
			existingUser = &user{id: string(id)}
			return nil
		})
		return existingUser
	}
}*/
