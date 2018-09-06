package commands

import (
	"github.com/bwmarrin/discordgo"
)

type CommandResultHandler interface {
	Handle(session *discordgo.Session, commandMessage *discordgo.MessageCreate, resultMessage CommandResultMessage)
}
