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

	/*channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}*/

	input := m.Content

	if m.Content[:1] == "<" && len(m.Content) >= 2 { // Called by mention
		mentionSize := len(s.State.User.ID) + 2
		idPrefix := 2
		if m.Content[2:3] == "!" {
			mentionSize++
			idPrefix++
		}
		if len(m.Content) < mentionSize+1 || m.Content[idPrefix:mentionSize] != s.State.User.ID {
			return
		}
		input = input[mentionSize+1:]
		m.Content = input
	} else {
		return
	}
	/* else if prefix, ok := prefixes[channel.GuildID]; ok && m.Content[:1] == prefix { // Called by prefix
		input = input[1:len(input)]
		h.Content = input
	}*/

	b.Parse(m)
}
