package commands

import "github.com/bwmarrin/discordgo"

// CommandError is an object satisfying the commandResultMessage interface
// Returned from a command handler, it signals the command processor that the request failed
type CommandError struct {
	CommandMessage *discordgo.MessageCreate
	Message        string
	Color          int
}

func (c *CommandError) setCommandMessage(commandMessage *discordgo.MessageCreate) {
	c.CommandMessage = commandMessage
}

func (c *CommandError) GetCommandMessage() *discordgo.MessageCreate {
	return c.CommandMessage
}

func (c *CommandError) GetMessage() string {
	return c.Message
}

func (c *CommandError) GetColor() int {
	return c.Color
}

func (c *CommandError) Error() string {
	return c.Message
}
