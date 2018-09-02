package db

import (
	bolt "go.etcd.io/bbolt"
)

// DB is a global reference to the database
var DB *bolt.DB

// InitDB connects the database to the global DB value
func Init() {
	var err error
	DB, err = bolt.Open("doddy.db", 0666, nil)
	if err != nil {
		panic("could not open boltdb: " + err.Error())
	}
}
