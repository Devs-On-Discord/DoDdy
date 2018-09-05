package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/bwmarrin/discordgo"
)

type botCommandValidator struct {
}

func (v botCommandValidator) Validate(command *commands.Command, session *discordgo.Session, message *discordgo.MessageCreate) (bool) {

	return true
}
