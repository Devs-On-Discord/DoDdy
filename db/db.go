package db

import (
	bolt "go.etcd.io/bbolt"
)

var Db *bolt.DB

func InitDb() {
	var err error
	Db, err = bolt.Open("doddy.db", 0666, nil)
	if err != nil {
		panic("could not open boltdb: " + err.Error())
	}
}
