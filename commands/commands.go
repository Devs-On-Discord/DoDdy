package commands

import (
	"github.com/anmitsu/go-shlex"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	commands         map[string]Command
	ResultMessages   chan commandResultMessage
	incomingMessages chan *discordgo.MessageCreate
}

func (c *Commands) Init() {
	c.commands = make(map[string]Command)
	c.ResultMessages = make(chan commandResultMessage)
	go func() {
		for {
			incomingMessage := <-c.incomingMessages
			c.parse(incomingMessage)
		}
	}()
}

func (c *Commands) Register(command Command) {
	commandNameSplit := strings.Split(command.Name, " ")
	if len(commandNameSplit) < 1 {
		return
	}
	name := commandNameSplit[0]
	c.commands[name] = command
}

func (c *Commands) parse(commandMessage *discordgo.MessageCreate) {
	commandParsed, err := shlex.Split(commandMessage.Content, true)
	if err != nil {
		c.ResultMessages <- commandError{
			commandMessage: commandMessage,
			message:        "Error happened " + err.Error(),
			color:          0xb30000,
		}
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		c.ResultMessages <- commandError{
			commandMessage: commandMessage,
			message:        "Invalid Command",
			color:          0xb30000,
		}
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[commandName]; exists {
		if commandCount < 2 {
			c.ResultMessages <- command.Handler(commandMessage, nil)
		} else {
			c.ResultMessages <- command.Handler(commandMessage, commandParsed[1:])
		}
	}
}

func (c *Commands) Parse(commandMessage *discordgo.MessageCreate) {
	c.incomingMessages <- commandMessage
}
