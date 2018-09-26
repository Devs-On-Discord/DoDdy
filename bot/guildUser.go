package main

import (
	"encoding/binary"
)

import (
	bolt "go.etcd.io/bbolt"
)

type guildUser struct {
	entity
	reputation uint64
	warns      []guildUserWarn
}

func (u *guildUser) Init() {
	u.entity.Init()
	u.fields = map[string]*entityField{
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
	u.name = "guildUser"
	u.onLoad = u.OnLoad
	u.onSave = u.OnSave
}

func (u *guildUser) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "reputation":
		return binary.LittleEndian.Uint64(val)
	}
	return nil
}

func (u *guildUser) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	return nil, nil
}
