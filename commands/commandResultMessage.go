package commands

import "github.com/bwmarrin/discordgo"

type commandResultMessage interface {
	CommandMessage() *discordgo.MessageCreate
	Session() *discordgo.Session
	Message() (string)
	Color() (int)
}
