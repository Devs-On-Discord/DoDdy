package main

import (
	"strconv"
)

import (
	bolt "go.etcd.io/bbolt"
)

type guildUser struct {
	entity
	reputation uint64
	warns      map[string]*guildUserWarn
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
		"warns": {
			getter: func() interface{} {
				return u.warns
			},
		},
	}
	u.name = "guildUser"
	u.onLoad = u.OnLoad
	u.onSave = u.OnSave

	u.warns = map[string]*guildUserWarn{}
}

func (u *guildUser) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "reputation":
		if reputation, err := strconv.ParseUint(string(val), 10, 64); err == nil {
			return reputation
		} else {
			return nil
		}
	case "warns":
		u.loadNestedBucketEntity(key, bucket, func(id string, bucket *bolt.Bucket) {
			guildUserWarn := &guildUserWarn{}
			guildUserWarn.Init()
			guildUserWarn.SetID(id)
			guildUserWarn.LoadBucket(bucket)
			u.warns[id] = guildUserWarn
		})
	}
	return nil
}

func (u *guildUser) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	switch key {
	case "warns":
		u.saveNestedBucketEntities(key, bucket, len(u.warns), func(save func(entity Entity)) {
			for _, warn := range u.warns {
				save(warn)
			}
		})
	}
	return nil, nil
}
