package commands

import "github.com/bwmarrin/discordgo"

type CommandIdentifier interface {
	Identify(s *discordgo.Session, m *discordgo.MessageCreate) bool
}
