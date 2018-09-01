package botcommands

import "github.com/bwmarrin/discordgo"
import "github.com/Devs-On-Discord/DoDdy/commands"

func helpCommand(commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	return &commands.CommandReply{Message: "Help command executed", Color: 0x00b300}
}
