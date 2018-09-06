package main

import bolt "go.etcd.io/bbolt"

type vote struct {
	entity
}

func (v *vote) Init() {
	v.entity.Init()
	v.name = "vote"
	v.onLoad = v.OnLoad
}

func (v *vote) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "name", "message":
		return string(val)
	case "answers":
		return v.loadNestedBucketEntityMap(key, bucket, func() Entity {
			return &voteAnswer{}
		})
	case "guild":
		return v.loadNestedBucketEntityMap(key, bucket, func() Entity {
			return &voteGuild{}
		})
	}
	return nil
}
