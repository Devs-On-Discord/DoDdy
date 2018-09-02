package votes

import "github.com/bwmarrin/discordgo"

type Votes struct {
	session *discordgo.Session
	votes   map[string]Vote
}

func (v *Votes) Init(session *discordgo.Session) {
	v.session = session
	v.session.AddHandler(v.reactionAdded)
	v.Load()
}

func (v *Votes) Load() {

}

func (v *Votes) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if vote, exists := v.votes[reaction.ChannelID]; exists {
		valid := false
		for _, answer := range vote.answers {
			if answer.emojiID == reaction.Emoji.ID {
				valid = true
			}
		}
		if !valid {
			go session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.ID, reaction.UserID)
		}
	}
}
