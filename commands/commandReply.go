package commands

import "github.com/bwmarrin/discordgo"

type CommandReply struct {
	CommandMessage *discordgo.MessageCreate
	Message string
	Color   int
}

func (c *CommandReply) setCommandMessage(commandMessage *discordgo.MessageCreate) {
	c.CommandMessage = commandMessage
}

func (c *CommandReply) commandMessage() *discordgo.MessageCreate {
	return c.CommandMessage
}

func (c *CommandReply) message() string {
	return c.Message
}

func (c *CommandReply) color() int {
	return c.Color
}
