package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// DiscordCommandResultHandler is a master object containing everything related to handling incoming commands
// as well as responding to them, and deleting answers 10 seconds later
type DiscordCommandResultHandler struct {
	commands            *Commands
	session             *discordgo.Session
}

// Init constructs the DiscordCommandResultHandler and launches the handling goroutines
func (d *DiscordCommandResultHandler) Init(commands *Commands, session *discordgo.Session) {
	d.commands = commands
	d.session = session
	go func() {
		for {
			commandResult := <-d.commands.ResultMessages
			switch commandResult.(type) {
			case *CommandReply:
				go func() {
					commandMessage := commandResult.commandMessage()
					message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
						Color: commandResult.color(),
						Title: commandResult.message(),
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Deletion in 10 seconds",
						},
					})
					time.Sleep(10 * time.Second)
					d.session.ChannelMessageDelete(message.ChannelID, message.ID)
					d.session.ChannelMessageDelete(commandMessage.ChannelID, commandMessage.ID)
				}()
			case *CommandError:
				go func() {
					commandMessage := commandResult.commandMessage()
					message, _ := d.session.ChannelMessageSendEmbed(commandMessage.ChannelID, &discordgo.MessageEmbed{
						Color: commandResult.color(),
						Title: commandResult.message(),
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Deletion in 10 seconds",
						},
					})
					time.Sleep(10 * time.Second)
					d.session.ChannelMessageDelete(message.ChannelID, message.ID)
					d.session.ChannelMessageDelete(commandMessage.ChannelID, commandMessage.ID)
				}()
			}
		}
	}()
}
