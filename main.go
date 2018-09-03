package main

import (
	"fmt"
	"github.com/Devs-On-Discord/DoDdy/db"
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"os"
	"os/signal"
	"syscall"

	"github.com/Devs-On-Discord/DoDdy/botCommands"
	"github.com/Devs-On-Discord/DoDdy/votes"
	"github.com/bwmarrin/discordgo"
)

const version = "0.0.1"

func main() {
	fmt.Printf("DoDdy %s starting\n", version)

	dataBase := db.Init()

	defer dataBase.Close()

	g := &guilds.Guilds{}
	g.Init(dataBase)

	bot, err := discordgo.New("Bot " + testToken)
	if err != nil {
		panic(err.Error())
	}

	if bot.Open() != nil {
		panic("could not open bot: " + err.Error())
	}

	defer bot.Close()

	botcommands.Init(g, bot)

	votes.Init(bot)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
}
