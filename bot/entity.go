package main

import (
	"strconv"

	bolt "go.etcd.io/bbolt"
)

// Entity is an interface to the underlying entity object
type Entity interface {
	Init()
	Id() string
	SetId(id string)
	Name() string
	Data() map[string]interface{}
	Set(key string, val interface{})
	Get(key string) interface{}
	Update(keys []string)
	Load()
	Delete()
	SetOnLoad(func(key string, val []byte, bucket *bolt.Bucket) interface{})
	LoadBucket(bucket *bolt.Bucket)
}

type entity struct {
	id     string
	name   string
	data   map[string]interface{}
	onLoad func(key string, val []byte, bucket *bolt.Bucket) interface{}
}

func (e *entity) Init() {
	e.data = map[string]interface{}{}
}

func (e entity) ID() string {
	return e.id
}

func (e *entity) SetId(id string) {
	e.id = id
}

func (e entity) Name() string {
	return e.name
}

func (e entity) Data() map[string]interface{} {
	return e.data
}

func (e *entity) Set(key string, val interface{}) {
	d := e.data
	d[key] = val
}

func (e entity) Get(key string) interface{} {
	d := e.data
	if data, exists := d[key]; exists {
		return data
	}
	return nil
}

func (e entity) GetString(key string) string {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case string:
			return data.(string)
		}
	}
	return ""
}

func (e entity) GetInt(key string) int {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case int:
			return data.(int)
		}
	}
	return 0
}

func (e entity) GetEntity(key string) *entity {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case *entity:
			return data.(*entity)
		}
	}
	return nil
}

func (e entity) GetEntities(key string) []*entity {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case []*entity:
			return data.([]*entity)
		}
	}
	return nil
}

func (e entity) Update(keys []string) {
	DB.Update(func(tx *bolt.Tx) error {
		if entitiesBucket, err := tx.CreateBucketIfNotExists([]byte(e.name)); err == nil {
			if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(e.id)); err == nil {
				if keys != nil { // Only update the provided keys when keys are provided
					for _, key := range keys {
						e.dbSet(key, entityBucket)
					}
				} else {
					e.dbSetAll(entityBucket)
				}
			}
		}
		return nil
	})
}

func (e entity) dbSet(key string, bucket *bolt.Bucket) {
	value := e.Get(key)
	if value == nil {
		bucket.Delete([]byte(key))
	} else {
		switch value.(type) {
		case string:
			bucket.Put([]byte(key), []byte(value.(string)))
		case int:
			bucket.Put([]byte(key), []byte(strconv.Itoa(value.(int))))
		case *entity:
			entity := value.(*entity)
			if entityBucket, err := bucket.CreateBucketIfNotExists([]byte(entity.name)); err == nil {
				entity.dbSetAll(entityBucket)
			}
		case []*entity:
			for _, entity := range value.([]*entity) {
				if entityBucket, err := bucket.CreateBucketIfNotExists([]byte(entity.name)); err == nil {
					entity.dbSetAll(entityBucket)
				}
			}
		}
	}
}

func (e entity) dbSetAll(bucket *bolt.Bucket) {
	for key := range e.data {
		e.dbSet(key, bucket)
	}
}

func (e *entity) LoadBucket(bucket *bolt.Bucket) {
	entityCursor := bucket.Cursor()
	for k, v := entityCursor.First(); k != nil; k, v = entityCursor.Next() {
		key := string(k)
		e.data[key] = e.onLoad(key, v, bucket)
	}
}

func (e *entity) Load() {
	DB.View(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(e.name)); entitiesBucket != nil {
			if entityBucket := entitiesBucket.Bucket([]byte(e.id)); entityBucket != nil {
				e.LoadBucket(entityBucket)
			}
		}
		return nil
	})
}

func (e *entity) SetOnLoad(onLoad func(key string, val []byte, bucket *bolt.Bucket) interface{}) {
	e.onLoad = onLoad
}

func (e entity) Delete() {
	DB.Update(func(tx *bolt.Tx) error {
		entitiesBucket := tx.Bucket([]byte(e.name))
		if entitiesBucket == nil {
			return nil
		}
		entitiesBucket.DeleteBucket([]byte(e.id))
		return nil
	})
}
