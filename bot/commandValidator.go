package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

type commandValidator struct {
	guilds *Guilds
}

func (v commandValidator) Validate(command *commands.Command, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	guild, err := v.guilds.Guild(message.GuildID)
	if err != nil {
		return false
	}
	member, err := session.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		return false
	}
	commandRole := RoleInt[command.Role]
	for role, id := range guild.Roles {
		for _, memberRole := range member.Roles {
			if id == memberRole {
				if role >= commandRole {
					return true
				}
			}
		}
	}

	return true
}
