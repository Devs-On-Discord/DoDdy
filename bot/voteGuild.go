package main

import bolt "go.etcd.io/bbolt"

type voteGuild struct {
	entity
	messageID string
	channelID string
}

func (g *voteGuild) Init() {
	g.entity.Init()
	g.fields = map[string]*entityField{
		"messageID": {
			setter: func(val interface{}) {
				if messageID, ok := val.(string); ok {
					g.messageID = messageID
				}
			},
			getter: func() interface{} {
				return g.messageID
			},
		},
		"channelID": {
			setter: func(val interface{}) {
				if channelID, ok := val.(string); ok {
					g.channelID = channelID
				}
			},
			getter: func() interface{} {
				return g.channelID
			},
		},
	}
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
