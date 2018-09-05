package commands

import "github.com/bwmarrin/discordgo"

type CommandValidator interface {
	Validate(command *Command, session *discordgo.Session, message *discordgo.MessageCreate) (bool)
}
