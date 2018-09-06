package main

import (
	"github.com/bwmarrin/discordgo"
)

type votes2 struct {
	entityCache
	// Is there for faster calculations in reaction add, remove
	channelVotes map[string]*vote // Key: channelID
}

func (v *votes2) Init(session *discordgo.Session) {
	v.entityCache.Init()
	v.name = "vote"
	v.onCreate = v.CreateEntity
	v.Entities()
	v.channelVotes = map[string]*vote{}
	v.fillChannelVotes()
	session.AddHandler(v.reactionAdded)
	session.AddHandler(v.reactionRemoved)
}

func (v *votes2) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if _, exists := v.channelVotes[reaction.ChannelID]; exists {
		/*for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				//IncreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}*/
		go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}

func (v *votes2) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if _, exists := v.channelVotes[reaction.ChannelID]; exists {
		/*for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				//DecreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}*/
		go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}

func (v *votes2) CreateEntity() Entity {
	vote := &vote{}
	return vote
}

func (v *votes2) UpdateEntity(entityPtr *Entity) {
	entity := *entityPtr
	vote := entity.(*vote)
	v.fillChannelVotesForVote(vote)
}

func (v *votes2) fillChannelVotesForVote(vote *vote) {
	if guilds, err := vote.GetEntitiesMap("guild"); err == nil {
		for _, guild := range guilds {
			if channelID, err := guild.GetString("channelID"); err == nil {
				v.channelVotes[channelID] = vote
			}
		}
	} else {
		println("error", err.Error())
	}
}

func (v *votes2) fillChannelVotes() {
	for _, entityPtr := range v.entities {
		entity := *entityPtr
		vote := entity.(*vote)
		v.fillChannelVotesForVote(vote)
	}
}
