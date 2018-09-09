package main

import "github.com/bwmarrin/discordgo"

type guilds struct {
	entityCache
	session *discordgo.Session
	channelQuestions map[string]*question
}

func (g *guilds) Init(session *discordgo.Session) {
	g.channelQuestions = map[string]*question{}
	g.session = session
	g.session.AddHandler(g.guildCreate)
	g.session.AddHandler(g.reactionAdded)
	g.session.AddHandler(g.reactionRemoved)
	g.entityCache.Init()
	g.name = "guild"
	g.onCreate = g.CreateEntity
	g.onUpdate = g.UpdateEntity
	g.Entities()
}

func (g *guilds) CreateEntity() Entity {
	guild := &guild{}
	return guild
}

func (g *guilds) UpdateEntity(entityPtr *Entity) {
	entity := *entityPtr
	guild := entity.(*guild)
	g.fillChannelQuestionsForQuestion(guild)
}

func (g *guilds) fillChannelQuestionsForQuestion(guild *guild) {
	for _, question := range guild.questions {
		g.channelQuestions[question.channelID] = question
	}
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

func (g *guilds) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
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

func (g *guilds) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
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
