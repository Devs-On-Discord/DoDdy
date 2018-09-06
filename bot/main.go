package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Devs-On-Discord/DoDdy/bot/commands"
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

	g2 := &guilds2{}
	g2.Init()

	v := &votes{}
	v.Init(bot.session)

	botCommands := &commands.Commands{}
	botCommands.Init(bot.session)
	botCommands.Validator = commandValidator{guilds: g2}
	botCommands.Identifier = commandIdentifier{guilds: g2}
	botCommands.ResultHandler = commandResultHandler{}
	botCommands.RegisterGroup(guildAdminCommands{guilds: g, votes: v, guilds2: g2})
	botCommands.RegisterGroup(helpCommands{botCommands})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
}
