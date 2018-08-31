package commands

import "github.com/bwmarrin/discordgo"

type commandReply struct {
	commandMessage *discordgo.MessageCreate
	session *discordgo.Session
	message string
	color   int
}

func (c commandReply) CommandMessage() *discordgo.MessageCreate {
	return c.commandMessage
}

func (c commandReply) Session() *discordgo.Session {
	return c.session
}

func (c commandReply) Message() string {
	return c.message
}

func (c commandReply) Color() int {
	return c.color
}
