package commands

import (
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	commands         map[string]Command
	ResultMessages   chan CommandResultMessage
	incomingMessages chan *discordgo.MessageCreate
}

func (c *Commands) Init() {
	c.commands = make(map[string]Command)
	c.ResultMessages = make(chan CommandResultMessage)
	c.incomingMessages = make(chan *discordgo.MessageCreate)
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
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Error happened " + err.Error(),
			Color:          0xb30000,
		}
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Invalid Command",
			Color:          0xb30000,
		}
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[commandName]; exists {
		if commandCount < 2 {
			resultMessage := command.Handler(commandMessage, nil)
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		} else {
			resultMessage := command.Handler(commandMessage, commandParsed[1:])
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		}
	} else {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Command doesn't exists",
			Color:          0xb30000,
		}
	}
}

func (c *Commands) Parse(commandMessage *discordgo.MessageCreate) {
	c.incomingMessages <- commandMessage
}
