package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/bwmarrin/discordgo"
)

// BotCommands is an object that encapsulates both Commands and a result handler
type BotCommands struct {
	commands                    *commands.Commands
	discordCommandResultHandler *commands.DiscordCommandResultHandler
}

// Init constructs the BotCommands object
func (b *BotCommands) Init(session *discordgo.Session) {
	b.commands = &commands.Commands{}
	b.commands.Init()
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.commands, session)
	b.RegisterCommands()
}

// RegisterCommands registers commands with the Commands object
func (b *BotCommands) RegisterCommands() {
	b.commands.Register(commands.Command{Name: "help", Handler: helpCommand})
}

// Parse is the input sink for commands
func (b *BotCommands) Parse(message *discordgo.MessageCreate) {
	b.commands.Parse(message)
}

func (b *BotCommands) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	b.Parse(m)
}
