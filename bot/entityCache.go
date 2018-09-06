package main

//var _ Entity = (*entity)(nil)

import bolt "go.etcd.io/bbolt"

// EntityCache is an interface to the underlying entityCache object
type EntityCache interface {
	Init()
	Name() string
	Entity(id string) *entity
}

type entityCache struct {
	name     string
	entities map[string]*entity
}

func (c *entityCache) Init() {
	c.entities = map[string]*entity{}
}

func (c entityCache) Name() string {
	return c.name
}

func (c *entityCache) Entity(id string) *entity {
	if entity, exists := c.entities[id]; exists {
		return entity
	}
	entity := &entity{}
	entity.id = id
	entity.name = c.name
	entity.Load()
	c.entities[id] = entity
	return entity
}

func (c *entityCache) Entities() *entityCache {
	DB.View(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(c.name)); entitiesBucket != nil {
			entitiesCursor := entitiesBucket.Cursor()
			for k, _ := entitiesCursor.First(); k != nil; k, _ = entitiesCursor.Next() {
				id := string(k)
				entity := &entity{}
				entity.id = id
				entity.name = c.name
				entity.LoadBucket(entitiesBucket)
				c.entities[id] = entity
			}
		}
		return nil
	})
	return c
}
