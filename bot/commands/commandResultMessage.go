package commands

import "github.com/bwmarrin/discordgo"

// CommandResultMessage encloses both CommandError and CommandReply, it allows for commands to signal the successfulness of a query
type CommandResultMessage interface {
	setCommandMessage(commandMessage *discordgo.MessageCreate)
	commandMessage() *discordgo.MessageCreate
	message() string
	color() int
}
