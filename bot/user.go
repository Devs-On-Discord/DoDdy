package main

import (
	"encoding/binary"
	bolt "go.etcd.io/bbolt"
)

type user struct {
	entity
	id   string
	reputation uint64
}

func (u *user) Init() {
	u.entity.Init()
	u.fields = map[string]*entityField{
		"id": {
			setter: func(val interface{}) {
				if id, ok := val.(string); ok {
					u.id = id
				}
			},
			getter: func() interface{} {
				return u.id
			},
		},
		"reputation": {
			setter: func(val interface{}) {
				if reputation, ok := val.(uint64); ok {
					u.reputation = reputation
				}
			},
			getter: func() interface{} {
				return u.reputation
			},
		},
	}
	u.name = "user"
	u.onLoad = u.OnLoad
	u.onSave = u.OnSave

	u.reputation = 0
}

func (u *user) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "id":
		return string(val)
	case "reputation":
		return binary.LittleEndian.Uint64(val)
	}
	return nil
}

func (u *user) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	return nil, nil
}
