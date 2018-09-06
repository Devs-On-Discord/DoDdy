package main

type votes2 struct {
	entityCache
	// Is there for faster calculations in reaction add, remove
	channelVotes map[string]*vote // Key: channelID
}

func (v *votes2) Init() {
	v.entityCache.Init()
	v.name = "vote"
	v.channelVotes = map[string]*vote{}
	v.onCreate = v.CreateEntity
	v.Entities()
	v.fillChannelVotes()
}

func (v *votes2) CreateEntity() Entity {
	vote := &vote{}
	return vote
}

func (v *votes2) fillChannelVotesForVote(vote *vote) {
	if guilds, err := vote.GetEntitiesMap("guild"); err == nil {
		for _, guild := range guilds {
			if channelID, err := guild.GetString("channelID"); err == nil {
				v.channelVotes[channelID] = vote
			}
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
