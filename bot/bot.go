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
}

func (b *bot) Close() {
	if err := b.session.Close(); err != nil {
		println("could not open session: " + err.Error())
	}
}
