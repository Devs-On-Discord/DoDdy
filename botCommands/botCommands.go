package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/Devs-On-Discord/DoDdy/db"
	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// Instance is a globally accessible BotCommands object
var Instance BotCommands

// Init constructs the global BotCommands object
func Init(session *discordgo.Session) {
	Instance = BotCommands{}
	Instance.Init(session)
}

// BotCommands is an object that encapsulates both Commands and a result handler
type BotCommands struct {
	commands                    *commands.Commands
	discordCommandResultHandler *commands.DiscordCommandResultHandler
	Prefixes                    map[string]string
}

// Init constructs the BotCommands object
func (b *BotCommands) Init(session *discordgo.Session) {
	b.Prefixes = make(map[string]string)
	db.DB.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket([]byte("guilds"))
		if guildsBucket != nil {
			guildsBucket.ForEach(func(k, v []byte) error {
				guildBucket := guildsBucket.Bucket(k)
				if guildBucket != nil {
					prefix := guildBucket.Get([]byte("prefix"))
					if prefix != nil {
						b.Prefixes[string(k)] = string(prefix)
					}
				}
				return nil
			})
		}
		return nil
	})
	b.commands = &commands.Commands{}
	b.commands.Init(session)
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.commands, session)
	b.RegisterCommands()
	session.AddHandler(b.messageHandler)
}

// RegisterCommands registers commands with the Commands object
func (b *BotCommands) RegisterCommands() {
	b.commands.Register(commands.Command{Name: "help", Handler: helpCommand})
	b.commands.Register(commands.Command{Name: "clearAnnouncements", Handler: clearAnnouncements})
	b.commands.Register(commands.Command{Name: "setAnnouncementsChannel", Handler: setAnnouncementsChannel})
	b.commands.Register(commands.Command{Name: "announce announcement", Handler: postAnnouncement})
	b.commands.Register(commands.Command{Name: "postLastMessageAsAnnouncement", Handler: postLastMessageAsAnnouncement})
	b.commands.Register(commands.Command{Name: "setVotesChannel", Handler: setVotesChannel})
	b.commands.Register(commands.Command{Name: "vote", Handler: postVote})
	b.commands.Register(commands.Command{Name: "prefix", Handler: setPrefix})
	b.commands.Register(commands.Command{Name: "setup", Handler: setup})
}

// Parse is the input sink for commands
func (b *BotCommands) Parse(message *discordgo.MessageCreate) {
	b.commands.Parse(message)
}

func (b *BotCommands) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	} else if len(input) > 1 {
		if prefix, exists := b.Prefixes[m.GuildID]; exists {
			if input[:1] == prefix {
				m.Content = input[1:]
			} else {
				return
			}
		} else {
			return
		}
	} else {
		return
	}
	/* else if prefix, ok := prefixes[channel.GuildID]; ok && m.Content[:1] == prefix { // Called by prefix
		input = input[1:len(input)]
		h.Content = input
	}*/

	b.Parse(m)
}
