package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

type commandValidator struct {
	guilds *guilds
}

func (v commandValidator) Validate(command *commands.Command, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	member, err := session.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		return false
	}
	guildPtr, err := v.guilds.Entity(message.GuildID)
	if err != nil {
		//TODO: check if member role is an high one to accept the first commands without setting the roles, for example !setup
		return false
	}
	guild := *guildPtr
	if rawRoles, err := guild.Get("roles"); err == nil {
		roles := rawRoles.(map[Role]string)
		commandRole := RoleInt[command.Role]
		for role, id := range roles {
			for _, memberRole := range member.Roles {
				if id == memberRole {
					if role >= commandRole {
						return true
					}
				}
			}
		}
	}
	return true
}
