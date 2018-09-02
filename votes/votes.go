package votes

import "github.com/bwmarrin/discordgo"

type Votes struct {
	session *discordgo.Session
}

func (v *Votes) Init(session *discordgo.Session) {
	v.session = session
	v.session.AddHandler(reactionAdded)
}

func reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {

}
