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

type DiscordCommandResultHandler struct {
	commands            *Commands
	erasableMessages    []erasableMessage
	newErasableMessages chan erasableMessage
	session             *discordgo.Session
}

func (d *DiscordCommandResultHandler) Init(commands *Commands, session *discordgo.Session) {
	d.commands = commands
	d.erasableMessages = make([]erasableMessage, 0)
	d.newErasableMessages = make(chan erasableMessage)
	d.session = session
	go func() {
		for {
			commandResult := <-d.commands.ResultMessages
			switch commandResult.(type) {
			case *CommandReply:
				commandMessage := commandResult.commandMessage()
				message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color: commandResult.color(),
					Title: commandResult.message(),
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
			case *CommandError:
				commandMessage := commandResult.commandMessage()
				message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
					Color: commandResult.color(),
					Title: commandResult.message(),
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
