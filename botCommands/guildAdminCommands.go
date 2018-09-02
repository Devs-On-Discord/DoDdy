package botcommands

import "github.com/bwmarrin/discordgo"
import "github.com/Devs-On-Discord/DoDdy/commands"
import (
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"github.com/Devs-On-Discord/DoDdy/votes"
)

func setVotesChannel(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	channel, err := session.Channel(channelID)
	if err != nil {
		return &commands.CommandError{Message: "Vote channel couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	err = guilds.SetVotesChannel(channel.GuildID, channelID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Vote channel set to " + channel.Name, Color: 0x00b300}
}

func postVote(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 3 {
		return &commands.CommandError{Message: "Vote id, name and message are required", Color: 0xb30000}
	}
	voteID := args[0]
	voteName := args[1]
	voteMessage := args[2]
	channels, err := guilds.GetVotesChannels()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	votes.AddVote(voteID, voteName, voteMessage, make([]votes.Answer, 0))
	for _, channelID := range channels {
		go func(channelID string) {
			message, err := session.ChannelMessageSend(channelID, voteMessage)
			if err == nil {
				channel, err := session.Channel(channelID)
				if err == nil {
					err = guilds.AddVote(channel.GuildID, voteID, message.ID, channelID)
					if err != nil {
						println(err.Error())
					}
					votes.Instance.Votes[channelID] = votes.Vote{Id: voteID, Name: voteName, Message: voteMessage, Answers: make([]votes.Answer, 0)}
				}
			}
		}(channelID)
	}
	return &commands.CommandReply{Message: "Vote posted", Color: 0x00b300}
}

func setAnnouncementsChannel(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	channel, err := session.Channel(channelID)
	if err != nil {
		return &commands.CommandError{Message: "Announcement channel couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	err = guilds.SetAnnouncementsChannel(channel.GuildID, channelID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Announcement channel set to " + channel.Name, Color: 0x00b300}
}

func postAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 1 {
		return &commands.CommandError{Message: "Announcement message missing", Color: 0xb30000}
	}
	announcement := args[0]
	channels, err := guilds.GetAnnouncementChannels()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	for _, channelID := range channels {
		go session.ChannelMessageSend(channelID, announcement)
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func clearAnnouncements(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channels, err := guilds.GetAnnouncementChannels()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	for _, channelID := range channels {
		messages, err := session.ChannelMessages(channelID, 100, "", "", "")
		if err == nil {
			messageIDs := make([]string, len(messages))
			for i, message := range messages {
				messageIDs[i] = message.ID
			}
			session.ChannelMessagesBulkDelete(channelID, messageIDs)
		} else {
			println(err.Error())
		}
	}
	return &commands.CommandReply{Message: "Announcements cleared", Color: 0x00b300}
}

func postLastMessageAsAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	messages, err := session.ChannelMessages(channelID, 1, commandMessage.Message.ID, "", "")
	if err != nil || len(messages) < 1 {
		return &commands.CommandError{Message: "Message couldn't be find " + err.Error(), Color: 0xb30000}
	}
	message := messages[0]
	if message == nil {
		return &commands.CommandError{Message: "Message couldn't be find " + err.Error(), Color: 0xb30000}
	}
	session.ChannelMessageDelete(channelID, message.ID)
	announcement := message.Content
	channels, err := guilds.GetAnnouncementChannels()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	for _, channelID := range channels {
		session.ChannelMessageSend(channelID, announcement)
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func setup(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	channel, err := session.Channel(channelID)
	if err != nil {
		return &commands.CommandError{Message: "Server couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	guildID := channel.GuildID
	guild, err := session.Guild(guildID)
	if err != nil {
		return &commands.CommandError{Message: "Server couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	err = guilds.Create(guildID, guild.Name)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "setup", Color: 0x00b300}
}
