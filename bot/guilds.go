package main

import "github.com/bwmarrin/discordgo"

type guilds struct {
	entityCache
	session *discordgo.Session
}

func (g *guilds) Init(session *discordgo.Session) {
	g.session = session
	g.session.AddHandler(g.guildCreate)
	g.entityCache.Init()
	g.name = "guild"
	g.onCreate = g.CreateEntity
	g.Entities()
}

func (g *guilds) CreateEntity() Entity {
	guild := &guild{}
	return guild
}

func (g *guilds) Guild(id string) (*guild, error) {
	entityPtr, err := g.Entity(id)
	if err != nil {
		return nil, err
	}
	guild, ok := (*entityPtr).(*guild)
	if !ok {
		return nil, &entityNotFoundError{}
	}
	return guild, nil
}

func (g *guilds) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	guild := event.Guild
	member, err := s.GuildMember(guild.ID, s.State.User.ID)
	if err != nil {
		//TODO: find channel where bot has permissions to send error message
		return
	}
	for _, role := range member.Roles {
		for _, guildRole := range guild.Roles {
			if guildRole.ID == role {
				if guildRole.Permissions&discordgo.PermissionAddReactions != 0 {

				}
				if guildRole.Permissions&discordgo.PermissionSendMessages != 0 {

				}
				if guildRole.Permissions&discordgo.PermissionManageMessages != 0 {

				}
				if guildRole.Permissions&discordgo.PermissionBanMembers != 0 {

				}
				if guildRole.Permissions&discordgo.PermissionManageChannels != 0 {

				}
			}
		}
	}

	guildEntity, err := g.Guild(guild.ID)
	if err == nil || guildEntity != nil {
		return //Guild already setup, do nothing
	}
	switch err.(type) {
	case entityNotFoundError:
		//TODO: create guild entity instead of doing it in !setup
		channel, err := s.GuildChannelCreate(guild.ID, "doddy-setup", discordgo.ChannelTypeGuildText)
		if err != nil {
			//TODO: find channel where bot has permissions to send error message
			return
		}
		s.ChannelMessageSend(channel.ID, "Use this channel to setup the bot. Type !setup for more infos.")
	}
}
