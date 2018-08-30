package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func handleMessageCreate(s *discordgo.Session, h *discordgo.MessageCreate) {
	if h.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(h.ChannelID, "Hello!")
	result, err := store.Collection("Users").
		Doc(fmt.Sprint(time.Now().Format("20060102150405"))).
		Set(context.Background(), map[string]string{"message": h.Content})
	if err != nil {
		fmt.Println("Could not save message: " + err.Error())
	}
	fmt.Printf("Result:\n%v", result)

}
