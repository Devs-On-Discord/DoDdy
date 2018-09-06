package main

type guilds struct {
	entityCache
}

func (g *guilds) Init() {
	g.entityCache.Init()
	g.name = "guild"
	g.onCreate = g.CreateEntity
	g.Entities()
}

func (g *guilds) CreateEntity() Entity {
	guild := &guild{}
	return guild
}
