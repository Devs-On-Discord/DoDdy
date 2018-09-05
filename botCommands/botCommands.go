package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"github.com/Devs-On-Discord/DoDdy/roles"
	"github.com/Devs-On-Discord/DoDdy/votes"
	"github.com/bwmarrin/discordgo"
)

// BotCommands is an object that encapsulates both Commands and a result handler
type BotCommands struct {
	commands                    *commands.Commands
	discordCommandResultHandler *commands.DiscordCommandResultHandler
	guilds                      *guilds.Guilds
	votes                       *votes.Votes
}

// Init constructs the BotCommands object
func (b *BotCommands) Init(guilds *guilds.Guilds, votes *votes.Votes, session *discordgo.Session) {
	b.guilds = guilds
	b.votes = votes
	b.commands = &commands.Commands{}
	b.commands.Init(session)
	b.commands.Validator = botCommandValidator{guilds: guilds}
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.commands, session)
	b.RegisterCommands()
	session.AddHandler(b.messageHandler)
}

// RegisterCommands registers commands with the Commands object
func (b *BotCommands) RegisterCommands() {
	guildAdminCommands := guildAdminCommands{guilds: b.guilds, votes: b.votes}
	helpCommands := helpCommands{b.commands}
	b.commands.Register(commands.Command{
		Name:        "prefix",
		Description: "Changes / Displays the prefix.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setPrefix,
	})
	b.commands.Register(commands.Command{
		Name:        "setAnnouncementsChannel",
		Description: "Redefines this node's announcement channel.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setAnnouncementsChannel,
	})
	b.commands.Register(commands.Command{
		Name:        "announce announcement",
		Description: "Post an announcement in this node.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postAnnouncement,
	})
	b.commands.Register(commands.Command{
		Name:        "clearAnnouncements",
		Description: "Empties this node's announcement channel.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.clearAnnouncements,
	})
	b.commands.Register(commands.Command{
		Name:        "postLastMessageAsAnnouncement",
		Description: "Repost the last message sent in this channel as an announcement",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postLastMessageAsAnnouncement,
	})
	b.commands.Register(commands.Command{
		Name:        "setVotesChannel",
		Handler:     guildAdminCommands.setVotesChannel,
		Role:        roles.NodeMod,
		Description: "Redefines this node's voting channel.",
	})
	b.commands.Register(commands.Command{
		Name:        "survey vote",
		Description: "Starts a DoD-wide survey.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postVote,
	})
	b.commands.Register(commands.Command{
		Name:        "setup",
		Description: "Modifies basic configuration settings",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setup,
	})
	b.commands.Register(commands.Command{
		Name:        "role",
		Description: "Specify roles",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setRole,
	})
	b.commands.Register(commands.Command{
		Name:        "roles",
		Description: "Get roles",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.getRoles,
	})
	b.commands.Register(commands.Command{
		Name:        "help",
		Description: "Begins a vote in this node's voting channel.",
		Role:        roles.NodeMod,
		Handler:     helpCommands.helpCommand,
	})
}

// Parse is the input sink for commands
func (b *BotCommands) Parse(session *discordgo.Session, message *discordgo.MessageCreate) {
	b.commands.Parse(session, message)
}

func (b *BotCommands) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	botID := s.State.User.ID
	if m.Author.ID == botID {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	valid := false
	//TODO: support multiple mentions and append all mentions at the bot reply
	for _, mention := range m.Mentions {
		if mention.ID == botID {
			m.Content = m.Content[len(mention.ID)+3:] //<@{botID}>
			valid = true
			break
		}
	}
	if !valid {
		input := m.Content
		if len(input) > 1 {
			guild, err := b.guilds.Guild(m.GuildID)
			if err != nil {
				return
			}
			if guild.Prefix == input[:1] {
				m.Content = m.Content[1:]
			} else {
				return
			}
		} else {
			return
		}
	}
	b.Parse(s, m)
}
