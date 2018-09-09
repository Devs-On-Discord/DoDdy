package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

type questionCommands struct {
}

func (q questionCommands) Init(session *discordgo.Session) {
	session.AddHandler(q.reactionAdded)
}

func (q questionCommands) Commands() []*commands.Command {
	return []*commands.Command{
		{
			Name:        "ask",
			Description: "ask question.",
			Role:        int(User),
			Handler:     q.ask,
		},
	}
}

func (q *questionCommands) ask(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 1 {
		return &commands.CommandReply{Message: "Question is required", Color: 0xb30000}
	}
	question := args[0]
	channel, err := session.GuildChannelCreateComplex(commandMessage.GuildID, discordgo.GuildChannelCreateData{Name: question})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	_, err = session.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Title: question,
		Color: 9001204,
		Author: &discordgo.MessageEmbedAuthor{
			Name: commandMessage.Author.Username,
			IconURL:  commandMessage.Author.AvatarURL("50x50"),
		},
	})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Created", Color: 0x00b300}
}

func (q *questionCommands) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	//TODO: check if channel is an question channel
	//TODO: check if user that added the reaction is channel owner
	//TODO: add 24hour time until channel remove after reaction
	//TODO: close channel conversations
	if reaction.Emoji.APIName() == "✅" {
	}
}

func (q *questionCommands) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	//TODO: check if channel is an question channel
	//TODO: check if user that added the reaction is channel owner
	//TODO: make channel conversation open again
	//TODO: stop deletion timer
	if reaction.Emoji.APIName() == "✅" {
	}
}
