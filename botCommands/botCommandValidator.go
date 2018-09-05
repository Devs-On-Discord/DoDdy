package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"github.com/bwmarrin/discordgo"
)

type botCommandValidator struct {
	guilds *guilds.Guilds
}

func (v botCommandValidator) Validate(command *commands.Command, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if guild, err := v.guilds.Guild(message.GuildID); err != nil {
		return false
	} else {
		member, err := session.GuildMember(message.GuildID, message.Author.ID)
		if err != nil {
			return false
		}
		for role, id := range guild.Roles {
			for _, memberRole := range member.Roles {
				if id == memberRole {
					if role >= command.Role {
						return true
					}
				}
			}
		}
	}
	return true
}
