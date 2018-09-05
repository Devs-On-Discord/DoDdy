package main

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/Devs-On-Discord/DoDdy/roles"
)

func RegisterCommands(g *Guilds, v *Votes, botCommands *commands.Commands) {
	//TODO: command !nodes that lists all guilds and there online count, maybe its possible to just embed an guild like in an invite
	guildAdminCommands := guildAdminCommands{guilds: g, votes: v}
	helpCommands := helpCommands{botCommands}
	botCommands.Register(commands.Command{
		Name:        "prefix",
		Description: "Changes / Displays the prefix.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setPrefix,
	})
	botCommands.Register(commands.Command{
		Name:        "setAnnouncementsChannel",
		Description: "Redefines this node's announcement channel.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setAnnouncementsChannel,
	})
	botCommands.Register(commands.Command{
		Name:        "announce announcement",
		Description: "Post an announcement in this node.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postAnnouncement,
	})
	botCommands.Register(commands.Command{
		Name:        "clearAnnouncements",
		Description: "Empties this node's announcement channel.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.clearAnnouncements,
	})
	botCommands.Register(commands.Command{
		Name:        "postLastMessageAsAnnouncement",
		Description: "Repost the last message sent in this channel as an announcement",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postLastMessageAsAnnouncement,
	})
	botCommands.Register(commands.Command{
		Name:        "setVotesChannel",
		Handler:     guildAdminCommands.setVotesChannel,
		Role:        roles.NodeMod,
		Description: "Redefines this node's voting channel.",
	})
	botCommands.Register(commands.Command{
		Name:        "survey vote",
		Description: "Starts a DoD-wide survey.",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.postVote,
	})
	botCommands.Register(commands.Command{
		Name:        "setup",
		Description: "Modifies basic configuration settings",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setup,
	})
	botCommands.Register(commands.Command{
		Name:        "role",
		Description: "Specify roles",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.setRole,
	})
	botCommands.Register(commands.Command{
		Name:        "roles",
		Description: "Get roles",
		Role:        roles.NodeMod,
		Handler:     guildAdminCommands.getRoles,
	})
	botCommands.Register(commands.Command{
		Name:        "help",
		Description: "Begins a vote in this node's voting channel.",
		Role:        roles.NodeMod,
		Handler:     helpCommands.helpCommand,
	})
}
