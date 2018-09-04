package botcommands

import (
	"fmt"

	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/bwmarrin/discordgo"
)

type helpCommands struct {
	commands *map[string]commands.Command
}

func (h *helpCommands) helpCommand(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	var helpText string

	for _, j := range *h.commands {
		helpText = fmt.Sprintf("%s%s: %s\n", helpText, j.Name, j.Description)
	}

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
			Message: "Can't send help via DM " + err.Error(),
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
		Message: "Help sent via DM",
		Color:   0x00b300,
	}
}
