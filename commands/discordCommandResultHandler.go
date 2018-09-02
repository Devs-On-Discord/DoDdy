package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

//TODO: change delete time for different commands
type erasableMessage struct {
	commandID  string
	answerID   string
	channelID  string
	expireTime time.Time
}

// DiscordCommandResultHandler is a master object containing everything related to handling incoming commands
// as well as responding to them, and deleting answers 10 seconds later
type DiscordCommandResultHandler struct {
	commands            *Commands
	erasableMessages    []erasableMessage
	newErasableMessages chan erasableMessage
	session             *discordgo.Session
}

// Init constructs the DiscordCommandResultHandler and launches the handling goroutines
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
					commandID:  commandMessage.ID,
					answerID:   message.ID,
					channelID:  commandMessage.ChannelID,
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
					commandID:  commandMessage.ID,
					answerID:   message.ID,
					channelID:  commandMessage.ChannelID,
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
						d.session.ChannelMessageDelete(erasableMessage.channelID, erasableMessage.commandID)
						d.session.ChannelMessageDelete(erasableMessage.channelID, erasableMessage.answerID)
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
