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

type entityField struct {
	setter func(newValue interface{})
	getter func() interface{}
}

// Entity is an interface to the underlying entity object
type Entity interface {
	Init()
	ID() string
	SetID(id string)
	Name() string
	Data() map[string]interface{}
	Fields() map[string]*entityField
	Set(key string, val interface{})
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	Get(key string) (interface{}, error)
	Update(keys []string) error
	Load() error
	Delete() error
	SetOnLoad(func(key string, val []byte, bucket *bolt.Bucket) interface{})
	SetOnSave(func(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error))
	LoadBucket(bucket *bolt.Bucket)
	dbSetAll(bucket *bolt.Bucket) error
}

type entity struct {
	id     string
	name   string
	data   map[string]interface{}
	fields map[string]*entityField
	onLoad func(key string, val []byte, bucket *bolt.Bucket) interface{}
	onSave func(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error)
}

func (e *entity) Init() {
	e.data = map[string]interface{}{}
	e.onLoad = e.onDefaultLoad
	e.onSave = e.onDefaultSave
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

func (e entity) Fields() map[string]*entityField {
	return e.fields
}

func (e *entity) Set(key string, val interface{}) {
	d := e.data
	d[key] = val
}

func (e entity) Get(key string) (interface{}, error) {
	if data, exists := e.data[key]; exists {
		return data, nil
	}
	return nil, &entityDataNotFoundError{key}
}

func (e entity) GetString(key string) (string, error) {
	if data, err := e.Get(key); err == nil {
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
	if data, err := e.Get(key); err == nil {
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
	if data, err := e.Get(key); err == nil {
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
	if data, err := e.Get(key); err == nil {
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
	if data, err := e.Get(key); err == nil {
		switch data.(type) {
		case map[string]Entity:
			return data.(map[string]Entity), nil
		default:
			return nil, &entityDataWrongTypeError{
				wrongType:   reflect.TypeOf(map[string]Entity{}),
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
				if keys != nil { // Only update the provided fields when fields are provided
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
	var value interface{}
	if field, exists := e.fields[key]; exists {
		if field == nil || field.getter == nil {
			// When getter doesn't exists make it possible to set the values manually in onSave with an nil value
			return e.dbSetValue(key, nil, bucket, false)
		} else {
			value = field.getter()
		}
	} else {
		value, _ = e.Get(key)
	}
	if value == nil {
		return bucket.Delete([]byte(key))
	} else {
		e.dbSetValue(key, value, bucket, false)
	}
	return nil
}

func (e entity) dbSetValue(key string, value interface{}, bucket *bolt.Bucket, secondRun bool) error {
	if value != nil {
		switch value.(type) {
		case string:
			return bucket.Put([]byte(key), []byte(value.(string)))
		case int:
			return bucket.Put([]byte(key), []byte(strconv.Itoa(value.(int))))
		case Entity:
			entity := value.(Entity)
			if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
				if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(entity.ID())); err == nil {
					return entity.dbSetAll(entityBucket)
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
		case map[string]string:
			if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
				for key, value := range value.(map[string]string) {
					err = entitiesBucket.Put([]byte(key), []byte(value))
				}
				return err
			} else {
				return err
			}
		}
	}
	value, err := e.onSave(key, value, bucket)
	if err != nil {
		return err
	}
	// Maybe onSave has given us an value we can actually save
	if value != nil && !secondRun {
		return e.dbSetValue(key, value, bucket, true)
	}
	return nil
}

func (e entity) dbSetAll(bucket *bolt.Bucket) error {
	var err error
	for field := range e.fields {
		err = e.dbSet(field, bucket)
	}
	for key := range e.data {
		err = e.dbSet(key, bucket)
	}
	return err
}

func (e *entity) LoadBucket(bucket *bolt.Bucket) {
	entityCursor := bucket.Cursor()
	for k, v := entityCursor.First(); k != nil; k, v = entityCursor.Next() {
		key := string(k)
		loadedValue := e.onLoad(key, v, bucket)
		if loadedValue == nil {
			continue
		}
		fields := e.fields //TODO: init fields in init when all are migrated
		if fields != nil {
			if field, exists := fields[key]; exists {
				if field != nil && field.setter != nil {
					field.setter(loadedValue)
				}
				// When there is no setter we assume its handled in onLoad, so setting it to data isn't needed
				continue
			}
		}
		e.data[key] = loadedValue
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

func (e *entity) SetOnSave(onSave func(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error)) {
	e.onSave = onSave
}

func (e entity) Delete() error {
	return DB.Update(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(e.name)); entitiesBucket != nil {
			return entitiesBucket.DeleteBucket([]byte(e.id))
		}
		return &bucketNotFoundError{e.name}
	})
}

func (e entity) loadNestedBucketEntity(key string, bucket *bolt.Bucket, loadEntity func(id string, bucket *bolt.Bucket)) {
	if nestedEntitiesBucket := bucket.Bucket([]byte(key)); nestedEntitiesBucket != nil {
		nestedEntitiesCursor := nestedEntitiesBucket.Cursor()
		for k, _ := nestedEntitiesCursor.First(); k != nil; k, _ = nestedEntitiesCursor.Next() {
			if nestedEntityBucket := nestedEntitiesBucket.Bucket(k); nestedEntityBucket != nil {
				loadEntity(string(k), nestedEntityBucket)
			}
		}
	}
}

func (e entity)saveNestedBucketEntities(key string, bucket *bolt.Bucket, count int, getEntities func(save func(entity Entity))) error {
	if entitiesBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
		save := func(entity Entity) {
			if entityBucket, err := entitiesBucket.CreateBucketIfNotExists([]byte(entity.ID())); err == nil {
				err = entity.dbSetAll(entityBucket)
			}
		}
		getEntities(save)
		return err
	} else {
		return err
	}
}

func (e entity) loadNestedBucket(key string, bucket *bolt.Bucket, loadValue func(key string, value string)) {
	if nestedValuesBucket := bucket.Bucket([]byte(key)); nestedValuesBucket != nil {
		nestedValuesCursor := nestedValuesBucket.Cursor()
		for k, v := nestedValuesCursor.First(); k != nil; k, v = nestedValuesCursor.Next() {
			loadValue(string(k), string(v))
		}
	}
}

func (e entity) saveCustomMap(key string, bucket *bolt.Bucket, customMap map[interface{}]interface{}, transformer func(key interface{}, value interface{}) ([]byte, []byte)) error {
	if customMapBucket, err := bucket.CreateBucketIfNotExists([]byte(key)); err == nil {
		var err error
		for key, value := range customMap {
			transformedKey, transformedValue := transformer(key, value)
			if transformedKey == nil || transformedValue == nil {
				continue
			}
			err = customMapBucket.Put(transformedKey, transformedValue)
		}
		return err
	} else {
		return err
	}
}

func (e entity) onDefaultLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	return nil
}

func (e entity) onDefaultSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	return nil, nil
}
