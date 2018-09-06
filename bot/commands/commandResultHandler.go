package commands

import (
	"github.com/bwmarrin/discordgo"
)

// CommandResultHandler is a master object containing everything related to handling incoming commands
// as well as responding to them, and deleting answers 10 seconds later
type CommandResultHandler interface {
	Handle(session *discordgo.Session, commandMessage *discordgo.MessageCreate, resultMessage CommandResultMessage)
}
