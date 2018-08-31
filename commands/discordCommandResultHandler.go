package commands

import "github.com/bwmarrin/discordgo"

type discordCommandResultHandler struct {
	commands *Commands
}

func (d *discordCommandResultHandler) Init() {
	go func() {
		for {
			commandResult := <-d.commands.ResultMessages
			switch commandResult.(type) {
			case commandReply:
				commandMessage := commandResult.CommandMessage()
				commandResult.Session().ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color:  commandResult.Color(),
					Title:  commandResult.Message(),
					Footer: &discordgo.MessageEmbedFooter{Text: "Deletion in 10 seconds"},
				})
			case commandError:
				commandMessage := commandResult.CommandMessage()
				commandResult.Session().ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color: commandResult.Color(),
					Title: commandResult.Message(),
				})
			}
		}
	}()
}
