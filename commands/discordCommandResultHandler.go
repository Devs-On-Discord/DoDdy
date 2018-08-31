package commands

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type erasableMessage struct {
	commandId  string
	answerId   string
	channelId  string
	expireTime time.Time
}

type discordCommandResultHandler struct {
	commands         *Commands
	erasableMessages chan erasableMessage
	session          *discordgo.Session
}

func (d *discordCommandResultHandler) Init() {
	d.erasableMessages = make(chan erasableMessage)
	go func() {
		for {
			commandResult := <-d.commands.ResultMessages
			switch commandResult.(type) {
			case commandReply:
				commandMessage := commandResult.CommandMessage()
				message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color: commandResult.Color(),
					Title: commandResult.Message(),
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Deletion in 10 seconds",
					},
				})
				d.erasableMessages <- erasableMessage{
					commandId:    commandMessage.ID,
					answerId:     message.ID,
					channelId:    commandMessage.ChannelID,
					expireTime: time.Now().Add(10 * time.Second),
				}
			case commandError:
				commandMessage := commandResult.CommandMessage()
				message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color: commandResult.Color(),
					Title: commandResult.Message(),
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Deletion in 10 seconds",
					},
				})
				d.erasableMessages <- erasableMessage{
					commandId:    commandMessage.ID,
					answerId:     message.ID,
					channelId:    commandMessage.ChannelID,
					expireTime: time.Now().Add(10 * time.Second),
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case x, ok := <-d.erasableMessages:

				if time.Now().After(x.expireTime) {
					d.session.ChannelMessageDelete(x.channelId, x.commandId)
					d.session.ChannelMessageDelete(x.channelId, x.answerId)
				} else {
					if ok {
						d.erasableMessages <- x
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		}
	}()
}
