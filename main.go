package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	bot, err := discordgo.New("Bot " + testToken)
	if err != nil {
		panic(err.Error())
	}
	if bot.Open() != nil {
		panic("could not open bot: " + err.Error())
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()
}
