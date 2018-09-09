package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

//TODO: move reactions handlers to questions entity cache struct
//TODO: add cache for channelID, question
type questionCommands struct {
	guilds *guilds
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

//TODO: maybe just ignore the args here and use the commandMessage without the command name as an question so the syntax with "" isn't needed here
func (q *questionCommands) ask(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 1 {
		return &commands.CommandReply{Message: "Question is required", Color: 0xb30000}
	}
	guild, err := q.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	questionMessage := args[0]
	commandMessage.Content = questionMessage
	questionMessage, err = commandMessage.ContentWithMoreMentionsReplaced(session)
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	channel, err := session.GuildChannelCreateComplex(commandMessage.GuildID, discordgo.GuildChannelCreateData{Name: questionMessage})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	_, err = session.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Title: questionMessage,
		Color: 0xEE2C90/*9001204*/,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    commandMessage.Author.Username,
			IconURL: commandMessage.Author.AvatarURL("50x50"),
		},
		Timestamp: string(commandMessage.Timestamp),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "@ RxJava",
			IconURL: "https://i.imgur.com/YP32xgu.png",
		},
	})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	if guild.questions == nil {
		guild.questions = map[string]*question{}
	}
	question := &question{}
	question.Init()
	question.id = channel.ID
	question.message = questionMessage
	question.askingUser = commandMessage.Author.ID
	question.channelID = channel.ID
	guild.questions[channel.ID] =  question
	guild.Update([]string{"questions"})
	q.guilds.Update(guild)
	return &commands.CommandReply{Message: "Created", Color: 0x00b300}
}
