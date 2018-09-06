package main

import (
	"fmt"
	"reflect"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type entityDataNotFoundError struct {
	key string
}

func (b entityDataNotFoundError) Error() string {
	return fmt.Sprintf("Entity Data (%s) couldn't be found", b.key)
}

type entityDataWrongTypeError struct {
	wrongType   reflect.Type
	correctType reflect.Type
}

func (b entityDataWrongTypeError) Error() string {
	return fmt.Sprintf("Entity Data doesn't confirm to %s is %s instead", b.wrongType, b.correctType)
}

// Entity is an interface to the underlying entity object
type Entity interface {
	Init()
	ID() string
	SetID(id string)
	Name() string
	Data() map[string]interface{}
	Set(key string, val interface{})
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	Get(key string) (interface{}, error)
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

func (e entity) Get(key string) (interface{}, error) {
	d := e.data
	if data, exists := d[key]; exists {
		return data, nil
	}
	return nil, &entityDataNotFoundError{key}
}

func (e entity) GetString(key string) (string, error) {
	if data, err := e.Get(key); err != nil {
		switch data.(type) {
		case string:
			return data.(string), nil
		default:
			return "", &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf(""),
				correctType: reflect.TypeOf(data),
			}
		}
	} else {
		return "", err
	}
}

func (e entity) GetInt(key string) (int, error) {
	if data, err := e.Get(key); err != nil {
		switch data.(type) {
		case int:
			return data.(int), nil
		default:
			return 0, &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf(0),
				correctType: reflect.TypeOf(data),
			}
		}
	} else {
		return 0, err
	}
}

func (e entity) GetEntity(key string) (*Entity, error) {
	if data, err := e.Get(key); err != nil {
		switch data.(type) {
		case *Entity:
			return data.(*Entity), nil
		default:
			return nil, &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf(&entity{}),
				correctType: reflect.TypeOf(data),
			}
		}
	} else {
		return nil, err
	}
}

func (e entity) GetEntities(key string) ([]*Entity, error) {
	if data, err := e.Get(key); err != nil {
		switch data.(type) {
		case []*Entity:
			return data.([]*Entity), nil
		default:
			return nil, &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf([]*entity{}),
				correctType: reflect.TypeOf(data),
			}
		}
	} else {
		return nil, err
	}
}

func (e entity) GetEntitiesMap(key string) (map[string]Entity, error) {
	if data, err := e.Get(key); err != nil {
		switch data.(type) {
		case map[string]Entity:
			return data.(map[string]Entity), nil
		default:
			return nil, &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf(map[string]entity{}),
				correctType: reflect.TypeOf(data),
			}
		}
	} else {
		return nil, err
	}
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
	value, err := e.Get(key)
	if err != nil {
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
				return nil
			}
			return &bucketNotFoundError{e.name + "-" + e.id}
		}
		return &bucketNotFoundError{e.name}
	})
}

func (e *entity) SetOnLoad(onLoad func(key string, val []byte, bucket *bolt.Bucket) interface{}) {
	e.onLoad = onLoad
}

func (e entity) Delete() error {
	return DB.Update(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(e.name)); entitiesBucket != nil {
			return entitiesBucket.DeleteBucket([]byte(e.id))
		}
		return &bucketNotFoundError{e.name}
	})
}
