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
	bot.AddHandler(func(s *discordgo.Session, h *discordgo.MessageCreate) {
		if h.Author.ID == s.State.User.ID {
			return
		}
		s.ChannelMessageSend(h.ChannelID, "Hello!")
	})
	if bot.Open() != nil {
		panic("could not open bot: " + err.Error())
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()
}
