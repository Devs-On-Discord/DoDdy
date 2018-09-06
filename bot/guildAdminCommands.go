package main

import (
	"bytes"
	"sync"

	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
)

/*
user
node mod with group (androiddev) (add channels to bot features)
node admin with group (androiddev) (add channels to bot features)
hub admins (add channels to bot features)
hub mods (add channels to bot features)
bot devs (add / remove channels to bot features)

make each role generatable single
and make existing roles as node roles

bot on join crease setup channel "doddy-setup"
when channel exists throw error
*/

type guildAdminCommands struct {
	guilds *Guilds
	votes  *votes2
}

func (g guildAdminCommands) Commands() []*commands.Command {
	return []*commands.Command{
		{
			Name:        "prefix",
			Description: "Changes / Displays the prefix.",
			Role:        int(NodeMod),
			Handler:     g.setPrefix,
		},
		{
			Name:        "setAnnouncementsChannel",
			Description: "Redefines this node's announcement channel.",
			Role:        int(NodeMod),
			Handler:     g.setAnnouncementsChannel,
		},
		{
			Name:        "announce announcement",
			Description: "Post an announcement in this node.",
			Role:        int(NodeMod),
			Handler:     g.postAnnouncement,
		},
		{
			Name:        "clearAnnouncements",
			Description: "Empties this node's announcement channel.",
			Role:        int(NodeMod),
			Handler:     g.clearAnnouncements,
		},
		{
			Name:        "postLastMessageAsAnnouncement",
			Description: "Repost the last message sent in this channel as an announcement",
			Role:        int(NodeMod),
			Handler:     g.postLastMessageAsAnnouncement,
		},
		{
			Name:        "setVotesChannel",
			Description: "Redefines this node's voting channel.",
			Role:        int(NodeMod),
			Handler:     g.setVotesChannel,
		},
		{
			Name:        "setVotesChannel",
			Description: "Redefines this node's voting channel.",
			Role:        int(NodeMod),
			Handler:     g.setVotesChannel,
		},
		{
			Name:        "survey vote",
			Description: "Starts a DoD-wide survey.",
			Role:        int(NodeMod),
			Handler:     g.postVote,
		},
		{
			Name:        "setup",
			Description: "Modifies basic configuration settings",
			Role:        int(NodeMod),
			Handler:     g.setup,
		},
		{
			Name:        "role",
			Description: "Specify roles",
			Role:        int(NodeMod),
			Handler:     g.setRole,
		},
		{
			Name:        "roles",
			Description: "Get roles",
			Role:        int(NodeMod),
			Handler:     g.getRoles,
		},
	}
}

func (g *guildAdminCommands) getRoles(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	guildRoles, err := session.GuildRoles(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: "Error in fetching server roles " + err.Error(), Color: 0xb30000}
	}

	var buffer bytes.Buffer
	for name, id := range CommandRoleNames {
		if role, exists := guild.Roles[id]; exists {
			roleName := role
			for _, guildRole := range guildRoles {
				if guildRole.ID == role {
					roleName = guildRole.Name
					break
				}
			}
			buffer.WriteString("role: " + name + " " + roleName + "\n")
		} else {
			buffer.WriteString("role: " + name + " not set\n")
		}
	}
	return &commands.CommandReply{Message: buffer.String(), Color: 0x00b300}
}

func (g *guildAdminCommands) setRole(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(commandMessage.MentionRoles) != 1 {
		return &commands.CommandError{Message: "Needs an single role mention", Color: 0xb30000}
	}
	if len(args) < 1 {
		return &commands.CommandError{Message: "Needs role name and role mention", Color: 0xb30000}
	}
	roleName := args[0]
	roleID := commandMessage.MentionRoles[0]
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	err = guild.SetRole(roleName, roleID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}

	return &commands.CommandReply{Message: "Role set", Color: 0x00b300}
}

func (g *guildAdminCommands) setPrefix(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	if len(args) < 1 {
		return &commands.CommandReply{Message: "Prefix is " + guild.Prefix, Color: 0xb30000}
	}
	prefix := args[0]
	err = guild.SetPrefix(prefix)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Bot prefix set to " + prefix, Color: 0x00b300}
}

func (g *guildAdminCommands) setVotesChannel(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	channel, err := session.Channel(channelID)
	if err != nil {
		return &commands.CommandError{Message: "Vote channel couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	err = guild.SetVotesChannel(channelID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Vote channel set to " + channel.Name, Color: 0x00b300}
}

//TODO: only create vote when it got successfully posted on all discord servers
func (g *guildAdminCommands) postVote(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 3 {
		return &commands.CommandError{Message: "Vote id, name and message are required", Color: 0xb30000}
	}
	voteID := args[0]
	voteName := args[1]
	voteMessage := args[2]

	loadedGuilds := g.guilds.LoadGuilds().Guilds

	var wg sync.WaitGroup
	wg.Add(len(loadedGuilds))

	voteGuilds := map[string]Entity{}

	go func() {
		wg.Wait()
		//TODO: error handling, add CommandResultMessage chan to Handle params
		vote := &vote{}
		vote.Init()
		vote.id = voteID
		vote.Set("name", voteName)
		vote.Set("message", voteMessage)
		vote.Set("guild", voteGuilds)

		vote.Update(nil)
		g.votes.Update(vote)
	}()

	for _, guild := range loadedGuilds {
		go func(channelID string) {
			defer wg.Done()
			message, err := session.ChannelMessageSend(channelID, voteMessage)
			if err == nil {
				channel, err := session.Channel(channelID)
				if err == nil {
					voteGuild := &voteGuild{}
					voteGuild.Init()
					voteGuild.id = channel.GuildID
					voteGuild.Set("channelID", channelID)
					voteGuild.Set("messageID", message.ID)
					voteGuilds[channel.GuildID] = voteGuild
				}
			}
		}(guild.VotesChannelID)
	}
	return &commands.CommandReply{Message: "Vote posted", Color: 0x00b300}
}

func (g *guildAdminCommands) setAnnouncementsChannel(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	channel, err := session.Channel(channelID)
	if err != nil {
		return &commands.CommandError{Message: "Announcement channel couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	err = guild.SetAnnouncementsChannel(channelID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Announcement channel set to " + channel.Name, Color: 0x00b300}
}

func (g *guildAdminCommands) postAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 1 {
		return &commands.CommandError{Message: "Announcement message missing", Color: 0xb30000}
	}
	announcement := args[0]
	for _, guild := range g.guilds.LoadGuilds().Guilds {
		go session.ChannelMessageSend(guild.AnnouncementsChannelID, announcement)
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func (g *guildAdminCommands) clearAnnouncements(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	for _, guild := range g.guilds.LoadGuilds().Guilds {
		messages, err := session.ChannelMessages(guild.AnnouncementsChannelID, 100, "", "", "")
		if err == nil {
			messageIDs := make([]string, len(messages))
			for i, message := range messages {
				messageIDs[i] = message.ID
			}
			session.ChannelMessagesBulkDelete(guild.AnnouncementsChannelID, messageIDs)
		} else {
			println(err.Error())
		}
	}
	return &commands.CommandReply{Message: "Announcements cleared", Color: 0x00b300}
}

func (g *guildAdminCommands) postLastMessageAsAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	messages, err := session.ChannelMessages(channelID, 1, commandMessage.Message.ID, "", "")
	if err != nil || len(messages) < 1 {
		return &commands.CommandError{Message: "Message couldn't be found " + err.Error(), Color: 0xb30000}
	}
	message := messages[0]
	if message == nil {
		return &commands.CommandError{Message: "Message couldn't be found " + err.Error(), Color: 0xb30000}
	}
	session.ChannelMessageDelete(channelID, message.ID)
	announcement := message.Content
	for _, guild := range g.guilds.LoadGuilds().Guilds {
		go session.ChannelMessageSend(guild.AnnouncementsChannelID, announcement)
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func (g *guildAdminCommands) setup(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guild, err := session.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: "Server couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	err = g.guilds.Create(commandMessage.GuildID, guild.Name)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "setup", Color: 0x00b300}
}
