package main

import (
	"time"

	"github.com/Devs-On-Discord/DoDdy/embed"

	"github.com/bwmarrin/discordgo"
)

func handleMessageCreate(s *discordgo.Session, h *discordgo.MessageCreate) {
	if h.Author.ID == s.State.User.ID {
		return
	}
	if len(h.Content) < len(s.State.User.ID)+3 {
		return
	}
	if h.Content[2:len(s.State.User.ID)+2] == s.State.User.ID {
		errMsg, _ := s.ChannelMessageSendEmbed(h.ChannelID, embed.NewEmbed().SetColor(0xFF0000).SetTitle("Command not recognized").SetFooter("Deletion in 10 seconds").MessageEmbed)
		time.Sleep(10 * time.Second)
		s.ChannelMessageDelete(h.ChannelID, h.ID)
		s.ChannelMessageDelete(h.ChannelID, errMsg.ID)
	}
}
