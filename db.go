package main

import (
	bolt "go.etcd.io/bbolt"
)

type db struct {
	db *bolt.DB
}

func (db *db) Init() {
	var err error
	db.db, err = bolt.Open("doddy.db", 0666, nil)
	if err != nil {
		panic("could not open boltdb: " + err.Error())
	}
}

func (db *db) Close() {
	db.db.Close()
}
