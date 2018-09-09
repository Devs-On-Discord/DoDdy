package commands

import "github.com/bwmarrin/discordgo"

// CommandReply is an object satisfying the commandResultMessage interface
// Returned from a command handler, it signals the command processor that the request was successful
type CommandReply struct {
	CommandMessage *discordgo.MessageCreate
	Message        string
	CustomMessage  *discordgo.MessageSend
	Color          int
}

func (c *CommandReply) setCommandMessage(commandMessage *discordgo.MessageCreate) {
	c.CommandMessage = commandMessage
}

func (c *CommandReply) GetCommandMessage() *discordgo.MessageCreate {
	return c.CommandMessage
}

func (c *CommandReply) GetMessage() string {
	return c.Message
}

func (c *CommandReply) GetCustomMessage() *discordgo.MessageSend {
	return c.CustomMessage
}

func (c *CommandReply) GetColor() int {
	return c.Color
}
