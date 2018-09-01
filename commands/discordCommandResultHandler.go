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
	commands            *Commands
	erasableMessages    []erasableMessage
	newErasableMessages chan erasableMessage
	session             *discordgo.Session
}

func (d *discordCommandResultHandler) Init() {
	d.erasableMessages = make([]erasableMessage, 0)
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
				d.newErasableMessages <- erasableMessage{
					commandId:  commandMessage.ID,
					answerId:   message.ID,
					channelId:  commandMessage.ChannelID,
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
				d.newErasableMessages <- erasableMessage{
					commandId:  commandMessage.ID,
					answerId:   message.ID,
					channelId:  commandMessage.ChannelID,
					expireTime: time.Now().Add(10 * time.Second),
				}
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				for i := len(d.erasableMessages) - 1; i >= 0; i-- {
					erasableMessage := d.erasableMessages[i]
					if time.Now().After(erasableMessage.expireTime) {
						d.session.ChannelMessageDelete(erasableMessage.channelId, erasableMessage.commandId)
						d.session.ChannelMessageDelete(erasableMessage.channelId, erasableMessage.answerId)
						d.erasableMessages = append(d.erasableMessages[:i], d.erasableMessages[i+1:]...)
					}
				}
				case newErasableMessage, ok := <-d.newErasableMessages:
					if ok {
						d.erasableMessages = append(d.erasableMessages, newErasableMessage)
					}
			}
		}
	}()
}
