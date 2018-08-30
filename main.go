package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/option"
)

const version = "0.0.1"

// The store is global for access in goroutines, this might create race conditions and lead to loss of data
// Miyoyo: I can't find anything saying this is goroutine safe or not, I'll assume, for the sake of simplicity
//         and because client is only a struct with data, that it is.
var store *firestore.Client

func main() {
	fmt.Printf("DoDdy %s starting\n", version)

	opt := option.WithCredentialsFile("firebase.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("could not initialize firebase: " + err.Error())
	}

	store, err = app.Firestore(context.Background())
	if err != nil {
		panic("could not connect to firestore: " + err.Error())
	}
	defer store.Close()

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
