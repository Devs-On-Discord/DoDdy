package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type commandIdentifier struct {
	guilds   *guilds
	replacer *strings.Replacer
}

func (i *commandIdentifier) Init(s *discordgo.Session) {
	i.replacer = strings.NewReplacer(
		"<@"+s.State.User.ID+">", "",
		"<@!"+s.State.User.ID+">", "",
	)
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
			m.Content = i.replacer.Replace(m.Content)
			//m.Content = m.Content[len(mention.ID)+3:] //<@{botID}>
			return true
		}
	}
	if len(m.Content) > 1 {
		if guildPtr, err := i.guilds.Entity(m.GuildID); err == nil {
			guild := (*guildPtr).(*guild)
			if guild.prefix == m.Content[:1] {
				m.Content = m.Content[1:]
				return true
			}
		}
	}
	return false
}
