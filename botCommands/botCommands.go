package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/bwmarrin/discordgo"
)

// BotCommands is an object that encapsulates both Commands and a result handler
type BotCommands struct {
	Commands                    *commands.Commands
	discordCommandResultHandler *commands.DiscordCommandResultHandler
}

// Init constructs the BotCommands object
func (b *BotCommands) Init(session *discordgo.Session) {
	b.Commands = &commands.Commands{}
	b.Commands.Init(session)
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.Commands, session)
	session.AddHandler(b.Commands.ProcessMessage)
}

func (b *BotCommands) Register(command commands.Command) {
	b.Commands.Register(command)
}
