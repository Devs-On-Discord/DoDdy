package main

import bolt "go.etcd.io/bbolt"

type question struct {
	entity
	askingUser string
	message string
	channelID string
}

func (q *question) Init() {
	q.entity.Init()
	q.fields = map[string]*entityField{
		"askingUser": {
			setter: func(val interface{}) {
				if askingUser, ok := val.(string); ok {
					q.askingUser = askingUser
				}
			},
			getter: func() interface{} {
				return q.askingUser
			},
		},
		"message": {
			setter: func(val interface{}) {
				if message, ok := val.(string); ok {
					q.message = message
				}
			},
			getter: func() interface{} {
				return q.message
			},
		},
		"channelID": {
			setter: func(val interface{}) {
				if channelID, ok := val.(string); ok {
					q.channelID = channelID
				}
			},
			getter: func() interface{} {
				return q.channelID
			},
		},
	}
	q.name = "question"
	q.onLoad = q.OnLoad
}

func (q *question) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "askingUser", "message", "channelID":
		return string(val)
	}
	return nil
}
