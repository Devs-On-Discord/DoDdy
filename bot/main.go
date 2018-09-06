package main

import (
	"fmt"
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"os"
	"os/signal"
	"syscall"
)

const version = "0.0.1"

func main() {
	fmt.Printf("DoDdy %s starting\n", version)

	db := db{}
	db.Init()

	defer db.Close()

	bot := bot{}
	bot.Init()

	defer bot.Close()

	g := &Guilds{}
	g.Init(db.db)

	v := &Votes{}
	v.Init(db.db, bot.session)

	botCommands := &commands.Commands{}
	botCommands.Init(bot.session)
	botCommands.Validator = commandValidator{guilds: g}
	botCommands.Identifier = commandIdentifier{guilds: g}
	botCommandsResultHandler := &commands.CommandResultHandler{}
	botCommandsResultHandler.Init(botCommands, bot.session)
	botCommands.RegisterGroup(guildAdminCommands{guilds: g, votes: v})
	botCommands.RegisterGroup(helpCommands{botCommands})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
}
