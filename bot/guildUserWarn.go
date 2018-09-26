package main

import (
	"encoding/binary"
)

import (
	bolt "go.etcd.io/bbolt"
)

type guildUserWarn struct {
	entity
	reason string
	authorID string
	timestamp uint64
}

func (u *guildUserWarn) Init() {
	u.entity.Init()
	u.fields = map[string]*entityField{
		"reason": {
			setter: func(val interface{}) {
				if reason, ok := val.(string); ok {
					u.reason = reason
				}
			},
			getter: func() interface{} {
				return u.reason
			},
		},
		"authorID": {
			setter: func(val interface{}) {
				if authorID, ok := val.(string); ok {
					u.authorID = authorID
				}
			},
			getter: func() interface{} {
				return u.authorID
			},
		},
		"timestamp": {
			setter: func(val interface{}) {
				if timestamp, ok := val.(uint64); ok {
					u.timestamp = timestamp
				}
			},
			getter: func() interface{} {
				return u.timestamp
			},
		},
	}
	u.name = "guildUserWarn"
	u.onLoad = u.OnLoad
	u.onSave = u.OnSave
}

func (u *guildUserWarn) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "reason", "authorID":
		return string(val)
	case "timestamp":
		return binary.LittleEndian.Uint64(val)
	}
	return nil
}

func (u *guildUserWarn) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	return nil, nil
}
