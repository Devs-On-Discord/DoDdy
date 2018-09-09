package main

import bolt "go.etcd.io/bbolt"

type voteAnswer struct {
	entity
	title     string
	emojiName string
	//count int //TODO: later
}

func (a *voteAnswer) Init() {
	a.entity.Init()
	a.fields = map[string]*entityField{
		"title": {
			setter: func(val interface{}) {
				if title, ok := val.(string); ok {
					a.title = title
				}
			},
			getter: func() interface{} {
				return a.title
			},
		},
		"emojiName": {
			setter: func(val interface{}) {
				if emojiName, ok := val.(string); ok {
					a.emojiName = emojiName
				}
			},
			getter: func() interface{} {
				return a.emojiName
			},
		},
	}
	a.name = "answer"
	a.onLoad = a.OnLoad
}

func (a *voteAnswer) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "title", "emojiName":
		return string(val)
	}
	return nil
}
