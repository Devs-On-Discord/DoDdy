package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

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
		{
			Name:        "topic",
			Description: "create topic.",
			Role:        int(NodeAdmin),
			Handler:     q.createTopic,
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
		Color: 0xEE2C90, /*9001204*/
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
	guild.questions[channel.ID] = question
	err = guild.Update([]string{"questions"})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	q.guilds.Update(guild)
	return &commands.CommandReply{Message: "Created", Color: 0x00b300}
}

func (q *questionCommands) createTopic(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 4 {
		return &commands.CommandReply{Message: "title, color, iconURL, child channel of topic", Color: 0xb30000}
	}
	title := args[0]
	color := args[1]
	iconURL := args[2]
	childChannel := args[3]
	childChannelID := childChannel[2 : len(childChannel)-1]

	guild, err := q.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}

	channel, err := session.Channel(childChannelID)
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}

	newTopic := &topic{}
	newTopic.Init()
	newTopic.id = channel.ParentID
	newTopic.title = title
	newTopic.color = "0x" + color
	newTopic.iconURL = iconURL
	newTopic.topicID = channel.ParentID
	newTopic.channelID = childChannelID
	if guild.topics == nil {
		guild.topics = map[string]*topic{}
	}
	guild.topics[channel.ParentID] = newTopic
	err = guild.Update([]string{"topics"})
	if err != nil {
		return &commands.CommandReply{Message: err.Error(), Color: 0xb30000}
	}
	q.guilds.Update(guild)

	return &commands.CommandReply{Message: "Created", Color: 0x00b300}
}
