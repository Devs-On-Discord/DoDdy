package main

import (
	"github.com/bwmarrin/discordgo"
)

type commandIdentifier struct {
	guilds *guilds
}

func (i commandIdentifier) Identify(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	botID := s.State.User.ID
	if m.Author.ID == botID {
		return false
	}
	if len(m.Content) == 0 {
		return false
	}
	for _, mention := range m.Mentions {
		if mention.ID == botID {
			m.Content = m.Content[len(mention.ID)+3:] //<@{botID}>
			return true
		}
	}
	if len(m.Content) > 1 {
		if guildPtr, err := i.guilds.Entity(m.GuildID); err == nil {
			guild := *guildPtr
			if prefix, err := guild.GetString("prefix"); err == nil {
				if prefix == m.Content[:1] {
					m.Content = m.Content[1:]
					return true
				}
			}
		}
	}
	return false
}
