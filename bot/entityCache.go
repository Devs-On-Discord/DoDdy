package main

//var _ Entity = (*entity)(nil)

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type entityNotFoundError struct {
}

func (b entityNotFoundError) Error() string {
	return fmt.Sprintf("Entity couldn't be found")
}

// EntityCache is an interface to the underlying entityCache object
type EntityCache interface {
	Init()
	Name() string
	Entity(id string) (*Entity, error)
}

type entityCache struct {
	name     string
	entities map[string]*Entity
	onCreate func() Entity
	onUpdate func(entity *Entity)
}

func (c *entityCache) Init() {
	c.entities = map[string]*Entity{}
	c.onCreate = c.createEntity
	c.onUpdate = c.updateEntity
}

func (c entityCache) Name() string {
	return c.name
}

func (c *entityCache) Entity(id string) (*Entity, error) {
	if entity, exists := c.entities[id]; exists {
		return entity, nil
	}
	entity := c.onCreate()
	entity.Init()
	entity.SetID(id)
	err := entity.Load()
	c.Update(entity)
	return &entity, err
}

func (c *entityCache) Entities() *entityCache {
	DB.View(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(c.name)); entitiesBucket != nil {
			entitiesCursor := entitiesBucket.Cursor()
			for k, _ := entitiesCursor.First(); k != nil; k, _ = entitiesCursor.Next() {
				if entityBucket := entitiesBucket.Bucket(k); entityBucket != nil {
					id := string(k)
					entity := c.onCreate()
					entity.Init()
					entity.SetID(id)
					entity.LoadBucket(entityBucket)
					c.Update(entity)
				}
			}
		}
		return nil
	})
	return c
}

func (c *entityCache) Update(entity Entity) {
	c.entities[entity.ID()] = &entity
	c.onUpdate(&entity)
}

func (c *entityCache) createEntity() Entity {
	return &entity{name: c.name}
}

func (c *entityCache) updateEntity(entity *Entity) {
}
