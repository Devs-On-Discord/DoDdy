package commands

import (
	"github.com/anmitsu/go-shlex"
	"strings"
)

type Commands struct {
	commands map[string]Command
}

func (c *Commands) Init() {
	c.commands = make(map[string]Command)
}

func (c *Commands) Register(command Command) {
	commandNameSplit := strings.Split(command.Name, " ")
	if len(commandNameSplit) < 1 {
		return
	}
	name := commandNameSplit[0]
	channel := make(chan []string)
	command.channel = channel
	c.commands[name] = command
	go func() {
		for {
			arguments := <-channel
			command.Handler(arguments)
		}
	}()
}

func (c *Commands) Parse(input string) (error) {
	commandParsed, err := shlex.Split(input, true)
	if err != nil {
		return err
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		return &commandError{message: "Invalid Command"}
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[commandName]; exists {
		if commandCount < 2 {
			command.channel <- nil
		} else {
			command.channel <- commandParsed[1:]
		}
	}
	return nil
}
