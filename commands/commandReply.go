package commands

import "github.com/bwmarrin/discordgo"

type commandReply struct {
	commandMessage *discordgo.MessageCreate
	message string
	color   int
}

func (c commandReply) CommandMessage() *discordgo.MessageCreate {
	return c.commandMessage
}

func (c commandReply) Message() string {
	return c.message
}

func (c commandReply) Color() int {
	return c.color
}
