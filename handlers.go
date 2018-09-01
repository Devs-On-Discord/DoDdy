package main

import (
	"fmt"
	"time"

	"github.com/Devs-On-Discord/DoDdy/embed"
	"github.com/anmitsu/go-shlex"
	"github.com/bwmarrin/discordgo"
)

var prefixes = map[string]string{}

func getPrefix(guildID string, username string) string {
	if prefix, ok := prefixes[guildID]; ok && prefixes[guildID] != "" {
		return prefix
	}
	return fmt.Sprintf("@%s ", username)
}

type deletionTarget struct {
	commandID    string
	answerID     string
	channelID    string
	deletionTime time.Time
}

// The message deletion channel is used to schedule messages for deletion, without having to keep a goroutine alive
var deletionChannel chan deletionTarget

func handleMessageCreate(s *discordgo.Session, h *discordgo.MessageCreate) {
	if h.Author.ID == s.State.User.ID {
		return
	}

	if len(h.Content) == 0 {
		return
	}

	channel, err := s.Channel(h.ChannelID)
	if err != nil {
		return
	}

	input := h.Content

	if h.Content[:1] == "<" && len(h.Content) >= 2 { // Called by mention
		mentionSize := len(s.State.User.ID) + 2
		idPrefix := 2
		if h.Content[2:3] == "!" {
			mentionSize++
			idPrefix++
		}
		if len(h.Content) < mentionSize+1 || h.Content[idPrefix:mentionSize] != s.State.User.ID {
			return
		}
		input = input[mentionSize+1 : len(input)]
		h.Content = input
	} else if prefix, ok := prefixes[channel.GuildID]; ok && h.Content[:1] == prefix { // Called by prefix
		input = input[1:len(input)]
		h.Content = input
	} else {
		return
	}

	commands.Parse(h)

	command, err := shlex.Split(input, true)
	if err != nil {
		answerThenDelete(h.ID, h.ChannelID, "could not parse:"+err.Error(), true, s)
		return
	}

	if len(command) == 0 || command[0] == "help" {
		dm, err := s.UserChannelCreate(h.Author.ID)
		if err != nil {
			answerThenDelete(h.ID, h.ChannelID, "Unable to initiate DM with the user.", true, s)
			return
		}
		_, err = s.ChannelMessageSendEmbed(dm.ID, embed.NewEmbed().SetTitle("Pretend this is the help string").MessageEmbed)
		if err != nil {
			answerThenDelete(h.ID, h.ChannelID, "Can't DM help, please allow DMs from this server.", true, s)
			return
		}
		s.ChannelMessageDelete(h.ChannelID, h.ID)
		return
	}

	answerThenDelete(h.ID, h.ChannelID, fmt.Sprintf("Command not recognized: %s", command[0]), true, s)
}

func answerThenDelete(commandID, channelID, message string, isError bool, s *discordgo.Session) {
	var color int
	if isError {
		color = 0xb30000
	} else {
		color = 0x00b300
	}
	answer, _ := s.ChannelMessageSendEmbed(channelID, embed.NewEmbed().SetColor(color).SetTitle(message).SetFooter("Deletion in 10 seconds").MessageEmbed)
	deletionChannel <- deletionTarget{
		commandID:    commandID,
		answerID:     answer.ID,
		channelID:    channelID,
		deletionTime: time.Now().Add(10 * time.Second),
	}
}

func deleter(input chan deletionTarget, s *discordgo.Session) {
	for {
		select {
		case x, ok := <-input:

			if time.Now().After(x.deletionTime) {
				s.ChannelMessageDelete(x.channelID, x.commandID)
				s.ChannelMessageDelete(x.channelID, x.answerID)
			} else {
				if ok {
					input <- x
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}
}
