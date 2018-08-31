package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const version = "0.0.1"

func main() {
	fmt.Printf("DoDdy %s starting\n", version)

	bot, err := discordgo.New("Bot " + testToken)
	if err != nil {
		panic(err.Error())
	}

	bot.AddHandler(handleMessageCreate)

	if bot.Open() != nil {
		panic("could not open bot: " + err.Error())
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
	bot.Close()
}
