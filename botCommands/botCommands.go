package botCommands

import (
	"../commands"
	"github.com/bwmarrin/discordgo"
)

type BotCommands struct {
	commands                    *commands.Commands
	session                     *discordgo.Session
	discordCommandResultHandler *commands.DiscordCommandResultHandler
}

func (b *BotCommands) Init() {
	b.commands = &commands.Commands{}
	b.commands.Init()
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.commands, b.session)
	b.RegisterCommands()
}

func (b *BotCommands) RegisterCommands() {
	b.commands.Register(commands.Command{Name: "help", Handler: helpCommand})
}

func (b *BotCommands) Parse(message *discordgo.MessageCreate) {
	b.commands.Parse(message)
}
