package main

import (
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type guild struct {
	entity
	serverName string
	prefix     string
	roles      map[Role]string
	channels   map[Channel]string
	questions  map[string]*question // Key: channelID
}

func (g *guild) Init() {
	g.entity.Init()
	g.fields = map[string]*entityField{
		"serverName": {
			setter: func(val interface{}) {
				if name, ok := val.(string); ok {
					g.serverName = name
				}
			},
			getter: func() interface{} {
				return g.serverName
			},
		},
		"prefix": {
			setter: func(val interface{}) {
				if prefix, ok := val.(string); ok {
					g.prefix = prefix
				}
			},
			getter: func() interface{} {
				return g.prefix
			},
		},
		"roles":     nil,
		"channels":  nil,
		"questions": nil,
	}
	g.name = "guild"
	g.onLoad = g.OnLoad
	g.onSave = g.OnSave
}

func (g *guild) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "serverName", "prefix":
		return string(val)
	case "roles":
		g.roles = map[Role]string{}
		g.loadNestedBucket(key, bucket, func(key string, value string) {
			if roleInt, err := strconv.Atoi(key); err == nil {
				if role, exists := RoleInt[roleInt]; exists {
					g.roles[role] = value
				}
			}
		})
	case "channels":
		g.channels = map[Channel]string{}
		g.loadNestedBucket(key, bucket, func(key string, value string) {
			if channelInt, err := strconv.Atoi(key); err == nil {
				if channel, exists := ChannelInt[channelInt]; exists {
					g.channels[channel] = value
				}
			}
		})
	case "questions":
		g.questions = map[string]*question{}
		g.loadNestedBucketEntity(key, bucket, func(id string, bucket *bolt.Bucket) {
			question := &question{}
			question.Init()
			question.SetID(id)
			question.LoadBucket(bucket)
			g.questions[id] = question
		})
	}

	return nil
}

func (g *guild) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	switch key {
	case "roles":
		if g.roles == nil {
			return nil, nil
		}
		if rolesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
			for role, roleId := range g.roles {
				err = rolesBucket.Put([]byte(strconv.Itoa(int(role))), []byte(roleId))
			}
			return nil, err
		} else {
			return nil, err
		}
	case "channels":
		if g.channels == nil {
			return nil, nil
		}
		if channelsBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
			for channel, channelId := range g.channels {
				err = channelsBucket.Put([]byte(strconv.Itoa(int(channel))), []byte(channelId))
			}
			return nil, err
		} else {
			return nil, err
		}
	case "questions":
		if g.questions == nil {
			return nil, nil
		}
		g.saveNestedBucketEntities(key, bucket, len(g.questions), func(save func(entity Entity)) {
			for _, question := range g.questions {
				save(question)
			}
		})
	}
	return nil, nil
}
