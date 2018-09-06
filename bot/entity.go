package main

import (
	"strconv"

	bolt "go.etcd.io/bbolt"
)

// Entity is an interface to the underlying entity object
type Entity interface {
	Init()
	ID() string
	SetID(id string)
	Name() string
	Data() map[string]interface{}
	Set(key string, val interface{})
	GetString(key string) string
	GetInt(key string) int
	Get(key string) interface{}
	Update(keys []string) error
	Load() error
	Delete() error
	SetOnLoad(func(key string, val []byte, bucket *bolt.Bucket) interface{})
	LoadBucket(bucket *bolt.Bucket)
	dbSetAll(bucket *bolt.Bucket) error
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

func (e *entity) SetID(id string) {
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

func (e entity) GetEntity(key string) *Entity {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case *Entity:
			return data.(*Entity)
		}
	}
	return nil
}

func (e entity) GetEntities(key string) []*Entity {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case []*Entity:
			return data.([]*Entity)
		}
	}
	return nil
}

func (e entity) GetEntitiesMap(key string) map[string]Entity {
	if data := e.Get(key); data != nil {
		switch data.(type) {
		case map[string]Entity:
			return data.(map[string]Entity)
		}
	}
	return nil
}

func (e entity) Update(keys []string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		if entitiesBucket, err := tx.CreateBucketIfNotExists([]byte(e.name)); err == nil {
			if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(e.id)); err == nil {
				if keys != nil { // Only update the provided keys when keys are provided
					for _, key := range keys {
						err = e.dbSet(key, entityBucket)
					}
					return err
				} else {
					return e.dbSetAll(entityBucket)
				}
			} else {
				return err
			}
		} else {
			return err
		}
		return nil
	})
}

func (e entity) dbSet(key string, bucket *bolt.Bucket) error {
	value := e.Get(key)
	if value == nil {
		return bucket.Delete([]byte(key))
	} else {
		switch value.(type) {
		case string:
			return bucket.Put([]byte(key), []byte(value.(string)))
		case int:
			return bucket.Put([]byte(key), []byte(strconv.Itoa(value.(int))))
		case Entity:
			entity := value.(Entity)
			if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
				if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(entity.ID())); err == nil {
					entity.dbSetAll(entityBucket)
				} else {
					return err
				}
			} else {
				return err
			}
		case []Entity:
			if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
				for _, entity := range value.([]Entity) {
					if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(entity.ID())); err == nil {
						err = entity.dbSetAll(entityBucket)
					}
				}
				return err
			} else {
				return err
			}
		case map[string]Entity:
			if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
				for key, entity := range value.(map[string]Entity) {
					if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(key)); err == nil {
						err = entity.dbSetAll(entityBucket)
					}
				}
				return err
			} else {
				return err
			}
		}
		return nil
	}
}

func (e entity) dbSetAll(bucket *bolt.Bucket) error {
	var err error
	for key := range e.data {
		err = e.dbSet(key, bucket)
	}
	return err
}

func (e *entity) LoadBucket(bucket *bolt.Bucket) {
	entityCursor := bucket.Cursor()
	for k, v := entityCursor.First(); k != nil; k, v = entityCursor.Next() {
		key := string(k)
		e.data[key] = e.onLoad(key, v, bucket)
	}
}

func (e *entity) Load() error {
	return DB.View(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(e.name)); entitiesBucket != nil {
			if entityBucket := entitiesBucket.Bucket([]byte(e.id)); entityBucket != nil {
				e.LoadBucket(entityBucket)
			} else {
				return &bucketNotFoundError{e.name + "-" + e.id}
			}
		} else {
			return &bucketNotFoundError{e.name}
		}
		return nil
	})
}

func (e *entity) SetOnLoad(onLoad func(key string, val []byte, bucket *bolt.Bucket) interface{}) {
	e.onLoad = onLoad
}

func (e entity) Delete() error {
	return DB.Update(func(tx *bolt.Tx) error {
		entitiesBucket := tx.Bucket([]byte(e.name))
		if entitiesBucket == nil {
			return &bucketNotFoundError{e.name}
		}
		err := entitiesBucket.DeleteBucket([]byte(e.id))
		return err
	})
}
