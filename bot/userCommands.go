package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

type userCommands struct {
}

func (u userCommands) Commands() []*commands.Command {
	return []*commands.Command{
		{
			Name:        "info",
			Description: "user infos.",
			Role:        int(User),
			Handler:     u.info,
		},
	}
}

//TODO: info without params is your own info
//TODO: anzeigen wenn gebannt
//TODO: warns anzeigen wenn man mod ist
//TODO: rank anzeigen
func (u *userCommands) info(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	user := user{}
	user.id = commandMessage.Author.ID
	user.Load()
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Reputation",
		Value:  strconv.FormatUint(user.reputation, 10),
		Inline: true,
	})
	return &commands.CommandReply{
		CustomMessage: &discordgo.MessageSend{
			Content: "user info",
			Embed: &discordgo.MessageEmbed{
				Color:  0x00b300, //TODO: use role color of the user
				Fields: fields,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    commandMessage.Author.Username,
					IconURL: commandMessage.Author.AvatarURL("50x50"),
				},
			},
		},
	}
}
