package main

import "github.com/bwmarrin/discordgo"

type bot struct {
	session *discordgo.Session
}

func (b *bot) Init() {
	var err error
	if b.session, err = discordgo.New("Bot " + testToken); err != nil {
		panic("could not initialize session: " + err.Error())
	}
	if err = b.session.Open(); err != nil {
		panic("could not open session: " + err.Error())
	}
	if b.session.State.User == nil {
		panic("No user")
	}
	b.session.AddHandler(b.ready)
}

func (b *bot) Close() {
	if err := b.session.Close(); err != nil {
		println("could not open session: " + err.Error())
	}
}

//TODO: fix ready event not getting called
func (b *bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	if err := s.UpdateStatus(0, "powering "+string(len(event.Guilds))+" nodes"); err != nil {
		println("update status error", err.Error())
	}
}
