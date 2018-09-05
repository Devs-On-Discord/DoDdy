package main

import bolt "go.etcd.io/bbolt"

type voteAnswer struct {
	entity
}

func (a *voteAnswer) Init() {
	a.entity.Init()
	a.name = "answer"
	a.onLoad = a.OnLoad
}

func (a *voteAnswer) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "name":
		return string(val)
	}
	return nil
}
