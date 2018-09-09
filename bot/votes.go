package main

import (
	"github.com/bwmarrin/discordgo"
)

type votes struct {
	entityCache
	// Is there for faster calculations in reaction add, remove
	channelVotes map[string]map[string]*voteGuild // Key: channelID, innerKey: messageID
}

func (v *votes) Init(session *discordgo.Session) {
	v.entityCache.Init()
	v.name = "vote"
	v.onCreate = v.CreateEntity
	v.onUpdate = v.UpdateEntity
	v.Entities()
	v.channelVotes = map[string]map[string]*voteGuild{}
	v.fillChannelVotes()
	session.AddHandler(v.reactionAdded)
	session.AddHandler(v.reactionRemoved)
}

func (v *votes) CreateEntity() Entity {
	vote := &vote{}
	return vote
}

func (v *votes) Vote(id string) (*vote, error) {
	entityPtr, err := v.Entity(id)
	if err != nil {
		return nil, err
	}
	vote, ok := (*entityPtr).(*vote)
	if !ok {
		return nil, &entityNotFoundError{}
	}
	return vote, nil
}

func (v *votes) UpdateEntity(entityPtr *Entity) {
	entity := *entityPtr
	vote := entity.(*vote)
	v.fillChannelVotesForVote(vote)
}

func (v *votes) fillChannelVotesForVote(vote *vote) {
	if vote.guild == nil {
		return
	}
	for _, guild := range vote.guild {
		if v.channelVotes[guild.channelID] == nil {
			v.channelVotes[guild.channelID] = map[string]*voteGuild{}
		}
		v.channelVotes[guild.channelID][guild.messageID] = guild
	}
}

func (v *votes) fillChannelVotes() {
	for _, entityPtr := range v.entities {
		vote := (*entityPtr).(*vote)
		v.fillChannelVotesForVote(vote)
	}
}

func (v *votes) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if guildVotes, exists := v.channelVotes[reaction.ChannelID]; exists {
		if _, exists := guildVotes[reaction.MessageID]; exists {
			/*for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				//IncreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}*/
			go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
			//go session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
		}
	}
}

func (v *votes) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
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
		//go session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
	}
}
