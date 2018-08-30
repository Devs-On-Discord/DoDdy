package main

import (
	"fmt"
	"time"

	"github.com/Devs-On-Discord/DoDdy/embed"
	"github.com/anmitsu/go-shlex"
	"github.com/bwmarrin/discordgo"
)

var prefixes = map[string]string{} // TODO: load prefixes from firebase on launch

func handleMessageCreate(s *discordgo.Session, h *discordgo.MessageCreate) {
	if h.Author.ID == s.State.User.ID {
		return
	}
	if len(h.Content) == 0 {
		return
	}

	input := h.Content
	channel, err := s.Channel(h.ChannelID)
	if err != nil {
		return
	}
	if h.Content[:1] == "<" && len(h.Content) >= 2 { // Called by mention
		nickSpacing := 0
		if h.Content[2:3] == "!" {
			nickSpacing = 1
		}
		if len(h.Content) >= len(s.State.User.ID)+3+nickSpacing && h.Content[2+nickSpacing:len(s.State.User.ID)+2+nickSpacing] == s.State.User.ID {
			input = input[len(s.State.User.ID)+3+nickSpacing : len(input)]
		} else {
			return
		}
	} else if prefix, ok := prefixes[channel.GuildID]; ok && h.Content[:1] == prefix { // Called by prefix
		input = input[1:len(input)]
	} else {
		return
	}

	command, err := shlex.Split(input, true)
	var problem = "No problem has been detected, please contact admins and report this."
	if err != nil {
		problem = "Could not parse command: " + err.Error()
	} else {
		fmt.Println(command)
		if len(command) == 0 || command[0] == "help" {
			dm, err := s.UserChannelCreate(h.Author.ID)
			if err != nil {
				problem = "Unable to initiate DM with the user."
			} else {
				_, err := s.ChannelMessageSendEmbed(dm.ID, embed.NewEmbed().SetTitle("Pretend this is the help string").MessageEmbed)
				if err != nil {
					problem = "Can't DM help, please allow DMs from this server."
				} else {
					s.ChannelMessageDelete(h.ChannelID, h.ID)
					return
				}
			}
		} else {
			problem = fmt.Sprintf("Command not recognized: %s", command[0])
		}
	}

	//Command routing should happen here

	errMsg, _ := s.ChannelMessageSendEmbed(h.ChannelID, embed.NewEmbed().SetColor(0xFF0000).SetTitle(problem).SetFooter("Deletion in 10 seconds").MessageEmbed)
	time.Sleep(10 * time.Second)
	s.ChannelMessageDelete(h.ChannelID, h.ID)
	s.ChannelMessageDelete(h.ChannelID, errMsg.ID)
}
