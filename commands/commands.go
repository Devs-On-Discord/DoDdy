package commands

import (
	"github.com/anmitsu/go-shlex"
	"strings"
)

type Commands struct {
	commands         map[string]Command
	ResultMessages   chan commandResultMessage
	incomingMessages chan string
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

func (c *Commands) parse(input string) {
	commandParsed, err := shlex.Split(input, true)
	if err != nil {
		c.ResultMessages <- commandError{message: "Error happened " + err.Error(), color: 0xb30000}
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		c.ResultMessages <- commandError{message: "Invalid Command", color: 0xb30000}
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[commandName]; exists {
		if commandCount < 2 {
			command.Handler(nil)
		} else {
			command.Handler(commandParsed[1:])
		}
	}
}

func (c *Commands) Parse(input string) {
	c.incomingMessages <- input
}
