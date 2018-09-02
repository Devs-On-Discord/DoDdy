package botcommands

import "github.com/bwmarrin/discordgo"
import "github.com/Devs-On-Discord/DoDdy/commands"

func helpCommand(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	helpText := "setup\n" +
		"setAnnouncementsChannel\n" +
		"postAnnouncement\n" +
		"postLastMessageAsAnnouncement\n" +
		"clearAnnouncements\n" +
		"setVotesChannel\n" +
		"vote\n" +
		"prefix"
	userChannel, err := session.UserChannelCreate(commandMessage.Author.ID)
	if err != nil {
		return &commands.CommandError{
			Message: "I couldn't contact you " + err.Error(),
			Color:   0xb30000,
		}
	}
	_, err = session.ChannelMessageSend(userChannel.ID, helpText)
	if err != nil {
		return &commands.CommandError{
			Message: "Help couldn't be send as an dm " + err.Error(),
			Color:   0xb30000,
		}
	}
	_, err = session.ChannelDelete(userChannel.ID)
	if err != nil {
		return &commands.CommandError{
			Message: "Couldn't cleanup channel " + err.Error(),
			Color:   0xb30000,
		}
	}
	return &commands.CommandReply{
		Message: "Help has been send as dm",
		Color:   0x00b300,
	}
}
