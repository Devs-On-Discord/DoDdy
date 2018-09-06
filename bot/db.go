package main

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type db struct {
	db *bolt.DB
}

type bucketNotFoundError struct {
	bucketName string
}

func (b bucketNotFoundError) Error() string {
	return fmt.Sprintf("Bucket (%s) couldn't be found", b.bucketName)
}

// DB Is a globally available database reference
var DB *bolt.DB

func (db *db) Init() {
	var err error
	db.db, err = bolt.Open("../doddy.db", 0666, nil)
	if err != nil {
		panic("could not open boltdb: " + err.Error())
	}
	DB = db.db
}

func (db *db) Close() {
	db.db.Close()
}
