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

func (g *guilds) Guild(id string) (*guild, error) {
	entityPtr, err := g.Entity(id)
	if err != nil {
		return nil, err
	}
	guild, ok := (*entityPtr).(*guild)
	if !ok {
		return nil, &entityNotFoundError{}
	}
	return guild, nil
}
