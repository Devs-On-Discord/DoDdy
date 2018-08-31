package commands

import "github.com/bwmarrin/discordgo"

type commandResultMessage interface {
	CommandMessage() *discordgo.MessageCreate
	Message() (string)
	Color() (int)
}
