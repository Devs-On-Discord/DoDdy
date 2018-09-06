package main

//var _ Entity = (*entity)(nil)

import bolt "go.etcd.io/bbolt"

type EntityCache interface {
	Init()
	Name() string
	Entity(id string) *Entity
}

type entityCache struct {
	name     string
	entities map[string]*Entity
	onCreate func() Entity
}

func (c *entityCache) Init() {
	c.entities = map[string]*Entity{}
	c.onCreate = c.createEntity
}

func (c entityCache) Name() string {
	return c.name
}

func (c *entityCache) Entity(id string) *Entity {
	if entity, exists := c.entities[id]; exists {
		return entity
	}
	entity := c.onCreate()
	entity.Init()
	entity.SetId(id)
	entity.Load()
	c.Update(entity)
	return &entity
}

func (c *entityCache) Entities() *entityCache {
	Db.View(func(tx *bolt.Tx) error {
		if entitiesBucket := tx.Bucket([]byte(c.name)); entitiesBucket != nil {
			entitiesCursor := entitiesBucket.Cursor()
			for k, _ := entitiesCursor.First(); k != nil; k, _ = entitiesCursor.Next() {
				if entityBucket := entitiesBucket.Bucket(k); entityBucket != nil {
					id := string(k)
					entity := c.onCreate()
					entity.Init()
					entity.SetId(id)
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
	c.entities[entity.Id()] = &entity
}

func (c *entityCache) createEntity() Entity {
	return &entity{name: c.name}
}
