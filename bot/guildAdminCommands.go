package main

import (
	"bytes"
	"sync"
	"time"

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
	votes  *votes
	guilds *guilds
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
		{
			Name:        "channel",
			Description: "Specify channels",
			Role:        int(NodeMod),
			Handler:     g.setChannel,
		},
		{
			Name:        "channels",
			Description: "Get channels",
			Role:        int(NodeMod),
			Handler:     g.getChannels,
		},
		{
			Name:        "warn",
			Description: "Warn user",
			Role:        int(NodeMod),
			Handler:     g.warn,
		},
	}
}

func (g *guildAdminCommands) getRoles(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guildRoles, err := session.GuildRoles(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: "Error in fetching server roles " + err.Error(), Color: 0xb30000}
	}

	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	var buffer bytes.Buffer
	for name, id := range CommandRoleNames {
		var role string
		exists := false
		if guild.channels != nil {
			role, exists = guild.roles[id]
		}
		if exists {
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

func (g *guildAdminCommands) setChannel(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 2 {
		return &commands.CommandError{Message: "Needs channel name and channel mention", Color: 0xb30000}
	}
	channelName := args[0]
	channelMention := args[1]
	channelID := channelMention[2 : len(channelMention)-1]
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	if commandChannel, exists := CommandChannelNames[channelName]; exists {
		guild.channels[commandChannel] = channelID
	} else {
		return &commands.CommandError{Message: "Unknown channel name " + channelName, Color: 0xb30000}
	}
	if err := guild.Update([]string{"channels"}); err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	g.guilds.Update(guild)
	return &commands.CommandReply{Message: "channel set", Color: 0x00b300}
}

func (g *guildAdminCommands) getChannels(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	var buffer bytes.Buffer
	for name, id := range CommandChannelNames {
		var channel string
		exists := false
		if guild.channels != nil {
			channel, exists = guild.channels[id]
		}
		if exists {
			buffer.WriteString("Channel: " + name + " <#" + channel + "> \n")
		} else {
			buffer.WriteString("channel: " + name + " not set\n")
		}
	}
	return &commands.CommandReply{
		CustomMessage: &discordgo.MessageSend{
			Content: buffer.String(),
			Embed: &discordgo.MessageEmbed{
				Color: 0x00b300,
				Title: "Deletion in 10 seconds",
			},
		},
	}
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
	if commandRole, exists := CommandRoleNames[roleName]; exists {
		guild.roles[commandRole] = roleID
	} else {
		return &commands.CommandError{Message: "Unknown role name " + roleName, Color: 0xb30000}
	}
	if err := guild.Update([]string{"roles"}); err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	g.guilds.Update(guild)
	return &commands.CommandReply{Message: "Role set", Color: 0x00b300}
}

func (g *guildAdminCommands) setPrefix(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	guild, err := g.guilds.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	if len(args) < 1 {
		return &commands.CommandReply{Message: "Prefix is " + guild.prefix, Color: 0xb30000}
	}
	prefix := args[0]
	guild.prefix = prefix
	err = guild.Update([]string{"prefix"})
	g.guilds.Update(guild)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	return &commands.CommandReply{Message: "Bot prefix set to " + prefix, Color: 0x00b300}
}

//TODO: only create vote when it got successfully posted on all discord servers
func (g *guildAdminCommands) postVote(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 3 {
		return &commands.CommandError{Message: "Vote id, name and message are required", Color: 0xb30000}
	}
	voteID := args[0]
	voteTitle := args[1]
	voteMessage := args[2]

	loadedEntities := g.guilds.Entities().entities

	var wg sync.WaitGroup
	wg.Add(len(loadedEntities))

	voteGuilds := map[string]*voteGuild{}

	go func() {
		wg.Wait()
		//TODO: error handling, add CommandResultMessage chan to Handle params
		vote := &vote{}
		vote.Init()
		vote.id = voteID
		vote.title = voteTitle
		vote.message = voteMessage
		vote.guild = voteGuilds
		if err := vote.Update(nil); err != nil {
			println("vote insert err", err.Error())
		}
		g.votes.Update(vote)
	}()

	for _, guildPtr := range loadedEntities {
		guild, ok := (*guildPtr).(*guild)
		if ok && guild.channels != nil {
			if channelID, exists := guild.channels[Votes]; exists {
				go func() {
					defer wg.Done()
					message, err := session.ChannelMessageSend(channelID, voteMessage)
					if err == nil {
						channel, err := session.Channel(channelID)
						if err == nil {
							voteGuild := &voteGuild{}
							voteGuild.Init()
							voteGuild.id = channel.GuildID
							voteGuild.channelID = channelID
							voteGuild.messageID = message.ID
							voteGuilds[channel.GuildID] = voteGuild
						}
					}
				}()
			} else {
				wg.Done()
			}
		} else {
			wg.Done()
		}
	}
	return &commands.CommandReply{Message: "Vote posted", Color: 0x00b300}
}

func (g *guildAdminCommands) postAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(args) < 1 {
		return &commands.CommandError{Message: "Announcement message missing", Color: 0xb30000}
	}
	announcement := args[0]
	for _, guildPtr := range g.guilds.Entities().entities {
		guild := (*guildPtr).(*guild)
		if channelID, exists := guild.channels[Announcements]; exists {
			go session.ChannelMessageSend(channelID, announcement)
		}
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func (g *guildAdminCommands) clearAnnouncements(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	for _, guildPtr := range g.guilds.Entities().entities {
		guild := (*guildPtr).(*guild)
		if channelID, exists := guild.channels[Announcements]; exists {
			messages, err := session.ChannelMessages(channelID, 100, "", "", "")
			if err == nil {
				messageIDs := make([]string, len(messages))
				for i, message := range messages {
					messageIDs[i] = message.ID
				}
				go session.ChannelMessagesBulkDelete(channelID, messageIDs)
			} else {
				println(err.Error())
			}
		}
	}
	return &commands.CommandReply{Message: "Announcements cleared", Color: 0x00b300}
}

func (g *guildAdminCommands) postLastMessageAsAnnouncement(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	channelID := commandMessage.ChannelID
	messages, err := session.ChannelMessages(channelID, 1, commandMessage.Message.ID, "", "")
	if err != nil || len(messages) < 1 || messages[0] == nil {
		return &commands.CommandError{Message: "Message couldn't be found " + err.Error(), Color: 0xb30000}
	}
	message := messages[0]
	session.ChannelMessageDelete(channelID, message.ID)
	announcement := message.Content
	for _, guildPtr := range g.guilds.Entities().entities {
		guild := (*guildPtr).(*guild)
		if channelID, exists := guild.channels[Announcements]; exists {
			go session.ChannelMessageSend(channelID, announcement)
		}
	}
	return &commands.CommandReply{Message: "Announcement posted", Color: 0x00b300}
}

func (g *guildAdminCommands) setup(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	sessionGuild, err := session.Guild(commandMessage.GuildID)
	if err != nil {
		return &commands.CommandError{Message: "Server couldn't be identified " + err.Error(), Color: 0xb30000}
	}
	newGuild := &guild{}
	newGuild.Init()
	newGuild.id = commandMessage.GuildID
	newGuild.serverName = sessionGuild.Name
	newGuild.prefix = "prefix"
	if err := newGuild.Update(nil); err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	g.guilds.Update(newGuild)

	return &commands.CommandReply{Message: "setup", Color: 0x00b300}
}

//TODO: pm user that he received an warn
func (g *guildAdminCommands) warn(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	if len(commandMessage.Mentions) != 1 {
		return &commands.CommandError{Message: "One User needs to be specified", Color: 0xb30000}
	}
	if len(args) < 2 {
		return &commands.CommandError{Message: "One reason needs to be specified", Color: 0xb30000}
	}
	reason := args[1]
	user := user{}
	user.Init()
	user.SetID(commandMessage.Mentions[0].ID)
	user.Load()

	warn := &guildUserWarn{}
	warn.Init()
	warn.reason = reason
	warn.authorID = commandMessage.Author.ID
	warn.timestamp = uint64(time.Now().Unix())

	if guild, ok := user.guilds[commandMessage.GuildID]; ok {
		warn.id = string(len(guild.warns) + 1)
		guild.warns[warn.id] = warn
		user.Update([]string{"guilds"})
	} else {
		warn.id = "1"
		guildUser := &guildUser{}
		guildUser.Init()
		guildUser.SetID(commandMessage.GuildID)
		guildUser.warns[warn.id] = warn
		user.guilds[commandMessage.GuildID] = guildUser
		user.Update(nil)

	}
	return &commands.CommandReply{Message: "User received the warn", Color: 0x00b300}
}
