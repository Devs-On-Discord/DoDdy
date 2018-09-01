package commands

import "github.com/bwmarrin/discordgo"

type CommandResultMessage interface {
	setCommandMessage(commandMessage *discordgo.MessageCreate)
	commandMessage() *discordgo.MessageCreate
	message() (string)
	color() (int)
}
