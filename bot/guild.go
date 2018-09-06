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
	}
	return nil
}

func (g *guild) OnSave(key string, val interface{}, bucket *bolt.Bucket) error {
	if key == "roles" {
		roles := val.(map[Role]string)
		if rolesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
			var err error
			for role, roleId := range roles {
				err = rolesBucket.Put([]byte(strconv.Itoa(int(role))), []byte(roleId))
			}
			return err
		} else {
			return err
		}
	}
	return nil
}
