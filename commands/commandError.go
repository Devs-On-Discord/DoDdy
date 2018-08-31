package commands

import "github.com/bwmarrin/discordgo"

type commandError struct {
	commandMessage *discordgo.MessageCreate
	message string
	color   int
}

func (c commandError) CommandMessage() *discordgo.MessageCreate {
	return c.commandMessage
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
