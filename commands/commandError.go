package commands

import "github.com/bwmarrin/discordgo"

type commandError struct {
	commandMessage *discordgo.MessageCreate
	session *discordgo.Session
	message string
	color   int
}

func (c commandError) CommandMessage() *discordgo.MessageCreate {
	return c.commandMessage
}

func (c commandError) Session() *discordgo.Session {
	return c.session
}

func (c commandError) Message() string {
	return c.message
}

func (c commandError) Color() int {
	return c.color
}

func (c *commandError) Error() string {
	return c.message
}
