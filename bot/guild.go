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
		"roles": nil,
		"channels": nil,
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
		if rolesBucket := bucket.Bucket([]byte(key)); rolesBucket != nil {
			g.roles = map[Role]string{}
			rolesCursor := rolesBucket.Cursor()
			for k, v := rolesCursor.First(); k != nil; k, v = rolesCursor.Next() {
				if roleInt, err := strconv.Atoi(string(k)); err == nil {
					if role, exists := RoleInt[roleInt]; exists {
						g.roles[role] = string(v)
					}
				}
			}
			return nil
		}
	case "channels":
		if channelsBucket := bucket.Bucket([]byte(key)); channelsBucket != nil {
			g.channels = map[Channel]string{}
			channelsCursor := channelsBucket.Cursor()
			for k, v := channelsCursor.First(); k != nil; k, v = channelsCursor.Next() {
				if channelInt, err := strconv.Atoi(string(k)); err == nil {
					if channel, exists := ChannelInt[channelInt]; exists {
						g.channels[channel] = string(v)
					}
				}
			}
			return nil
		}
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
	}
	return nil, nil
}
