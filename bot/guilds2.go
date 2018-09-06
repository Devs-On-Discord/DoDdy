package main

type guilds2 struct {
	entityCache
}

func (g *guilds2) Init() {
	g.entityCache.Init()
	g.name = "guild"
	g.onCreate = g.CreateEntity
	g.Entities()
}

func (g *guilds2) CreateEntity() Entity {
	guild := &guild{}
	return guild
}
