package main

import bolt "go.etcd.io/bbolt"

type voteGuild struct {
	entity
}

func (g *voteGuild) Init() {
	g.entity.Init()
	g.name = "guild"
	g.onLoad = g.OnLoad
}

func (g *voteGuild) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "messageID", "channelID":
		return string(val)
	}
	return nil
}
