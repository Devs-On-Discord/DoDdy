package main

type votes2 struct {
	entityCache
	// Is there for faster calculations in reaction add, remove
	channelVotes map[string]*vote // Key: channelID
}

func (v *votes2) Init() {
	v.entityCache.Init()
	v.name = "vote"
	v.onCreate = v.CreateEntity
	v.Entities()
	v.fillChannelVotes()
}

func (v *votes2) CreateEntity() Entity {
	vote := &vote{}
	return vote
}

func (v *votes2) fillChannelVotesForVote(vote *vote) {
	if guilds := vote.GetEntitiesMap("guild"); guilds != nil {
		for _, guild := range guilds {
			v.channelVotes[guild.GetString("channelID")] = vote
		}
	}
}

func (v *votes2) fillChannelVotes() {
	for _, entityPtr := range v.entities {
		entity := *entityPtr
		vote := entity.(*vote)
		v.fillChannelVotesForVote(vote)
	}
}
