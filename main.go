package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Devs-On-Discord/DoDdy/botCommands"
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

	botCommands := &botcommands.BotCommands{}
	botCommands.Init(bot.session)
	botCommands.Commands.Validator = commandValidator{guilds: g}
	botCommands.Commands.Identifier = commandIdentifier{guilds: g}

	RegisterCommands(g, v, botCommands)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
}
