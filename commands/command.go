package commands

import "github.com/bwmarrin/discordgo"


type Handle func(commandMessage *discordgo.MessageCreate, args []string) (commandResultMessage)

type Command struct {
	Name string
	Handler Handle
}
