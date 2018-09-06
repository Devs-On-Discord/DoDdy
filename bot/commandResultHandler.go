package main

import (
	"time"

	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

type commandResultHandler struct {
}

func (h commandResultHandler) Handle(session *discordgo.Session, commandMessage *discordgo.MessageCreate, resultMessage commands.CommandResultMessage) {
	switch resultMessage.(type) {
	case *commands.CommandReply, *commands.CommandError:
		go func() {
			message, err := session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
				Color: resultMessage.GetColor(),
				Title: resultMessage.GetMessage(),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Deletion in 10 seconds",
				},
			})
			if err == nil {
				time.Sleep(10 * time.Second)
				session.ChannelMessageDelete(message.ChannelID, message.ID)
				session.ChannelMessageDelete(commandMessage.ChannelID, commandMessage.ID)
			}
		}()
	}
}
