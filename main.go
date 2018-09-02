package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Devs-On-Discord/DoDdy/db"
	"github.com/bwmarrin/discordgo"
	"github.com/Devs-On-Discord/DoDdy/votes"
	"github.com/Devs-On-Discord/DoDdy/botCommands"
)

const version = "0.0.1"

func main() {
	fmt.Printf("DoDdy %s starting\n", version)

	bot, err := discordgo.New("Bot " + testToken)
	if err != nil {
		panic(err.Error())
	}

	db.InitDB()

	defer db.DB.Close()

	botcommands.Init(bot)

	votes.Init(bot)

	//TODO: Reimplement prefixes
	/*if db.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Nodes"))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				prefixes[string(k)] = string(b.Bucket(k).Get([]byte("Prefix")))
			}
			return nil
		}) != nil {
		panic("could not read prefixes from boltdb: " + err.Error())
	}*/

	if bot.Open() != nil {
		panic("could not open bot: " + err.Error())
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("DoDdy ready.\nPress Ctrl+C to exit.")
	<-sc
	bot.Close()
}
