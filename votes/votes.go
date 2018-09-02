package votes

import (
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"github.com/bwmarrin/discordgo"
)

var Instance Votes

func Init(session *discordgo.Session) {
	Instance = Votes{}
	Instance.Init(session)
}

type Votes struct {
	session *discordgo.Session
	Votes   map[string]Vote
}

func (v *Votes) Init(session *discordgo.Session) {
	v.session = session
	v.session.AddHandler(v.reactionAdded)
	v.session.AddHandler(v.reactionRemoved)
	v.Votes = make(map[string]Vote, 0)
	v.Load()
}

func (v *Votes) Load() {
	votes, err := GetVotes()
	if err != nil {
		println(err.Error())
	}
	guildVotes, err := guilds.GetVotes()
	if err != nil {
		println(err.Error())
	}
	var guildVoteVote *Vote = nil
	for _, guildVote := range guildVotes {
		for _, vote := range votes {
			if vote.Id == guildVote.VoteID {
				guildVoteVote = &vote
			}
		}
		if guildVoteVote != nil {
			v.Votes[guildVote.ChannelID] = *guildVoteVote
			guildVoteVote = nil
		}
	}
}

func (v *Votes) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if vote, exists := v.Votes[reaction.ChannelID]; exists {
		for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				IncreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}
		session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}

func (v *Votes) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if vote, exists := v.Votes[reaction.ChannelID]; exists {
		for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				DecreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}
		session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}
