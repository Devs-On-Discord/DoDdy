package main

import (
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type guild struct {
	entity
}

func (g *guild) Init() {
	g.entity.Init()
	g.name = "guild"
	g.onLoad = g.OnLoad
	g.onSave = g.OnSave
}

func (g *guild) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "name", "prefix":
		return string(val)
	case "roles":
		if rolesBucket := bucket.Bucket([]byte(key)); rolesBucket != nil {
			roles := map[Role]string{}
			rolesCursor := rolesBucket.Cursor()
			for k, v := rolesCursor.First(); k != nil; k, v = rolesCursor.Next() {
				if roleInt, err := strconv.Atoi(string(k)); err == nil {
					if role, exists := RoleInt[roleInt]; exists {
						roles[role] = string(v)
					}
				}
			}
			return roles
		}
	case "channels":
		if channelsBucket := bucket.Bucket([]byte(key)); channelsBucket != nil {
			channels := map[Channel]string{}
			channelsCursor := channelsBucket.Cursor()
			for k, v := channelsCursor.First(); k != nil; k, v = channelsCursor.Next() {
				if channelInt, err := strconv.Atoi(string(k)); err == nil {
					if channel, exists := ChannelInt[channelInt]; exists {
						channels[channel] = string(v)
					}
				}
			}
			return channels
		}
	}
	return nil
}

func (g *guild) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	switch key {
	case "roles":
		roles := val.(map[Role]string)
		if rolesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
			var err error
			for role, roleId := range roles {
				err = rolesBucket.Put([]byte(strconv.Itoa(int(role))), []byte(roleId))
			}
			return nil, err
		} else {
			return nil, err
		}
	case "channels":
		channels := val.(map[Channel]string)
		if channelsBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
			var err error
			for channel, channelId := range channels {
				err = channelsBucket.Put([]byte(strconv.Itoa(int(channel))), []byte(channelId))
			}
			return nil, err
		} else {
			return nil, err
		}
	}
	return nil, nil
}
