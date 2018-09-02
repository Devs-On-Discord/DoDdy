package botcommands

import "github.com/bwmarrin/discordgo"
import "github.com/Devs-On-Discord/DoDdy/commands"

func helpCommand(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	return &commands.CommandReply{Message: "!setup\n!setAnnouncementsChannel\n!postAnnouncement\n!postLastMessageAsAnnouncement\n!clearAnnouncements", Color: 0x00b300}
}
