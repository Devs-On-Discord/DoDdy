package main

import (
	bolt "go.etcd.io/bbolt"
)

type user struct {
	entity
	guilds map[string]*guildUser // Key: guildID
}

func (u *user) Init() {
	u.entity.Init()
	u.fields = map[string]*entityField{
	}
	u.name = "user"
	u.onLoad = u.OnLoad
	u.onSave = u.OnSave
}

func (u *user) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "guilds":
		u.guilds = map[string]*guildUser{}
		u.loadNestedBucketEntity(key, bucket, func(id string, bucket *bolt.Bucket) {
			guildUser := &guildUser{}
			guildUser.Init()
			guildUser.SetID(id)
			guildUser.LoadBucket(bucket)
			u.guilds[id] = guildUser
		})
	}
	return nil
}

func (u *user) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	switch key {
	case "answers":
		if u.guilds == nil {
			return nil, nil
		}
		u.saveNestedBucketEntities(key, bucket, len(u.guilds), func(save func(entity Entity)) {
			for _, guild := range u.guilds {
				save(guild)
			}
		})
	}
	return nil, nil
}
