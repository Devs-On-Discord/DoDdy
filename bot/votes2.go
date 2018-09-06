package main

type votes2 struct {
	entityCache
}

func (v *votes2) Init() {
	v.entityCache.Init()
	v.name = "vote"
	v.onCreate = v.CreateEntity
	v.Entities()
	//TODO: prepare channel cache
}

func (v *votes2) CreateEntity() Entity {
	vote := &vote{}
	return vote
}

/*func (v *votes2) OnCreate() *Entity {
	vote := vote2{}
	entity := Entity(vote)
	return &entity
}*/
