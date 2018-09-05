package commands

import (
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/bwmarrin/discordgo"
)

// Commands is an object containing all commands that can be called, it converts discordgo commands to a single thread each
// Commands contains registered commands ready to be called
// ResultMessages is the return channel for successful commands
type Commands struct {
	RegisteredCommands []*Command
	commands           map[string]Command
	ResultMessages     chan CommandResultMessage
	session            *discordgo.Session
	Validator          CommandValidator
}

// Init constructs the Commands object
func (c *Commands) Init(session *discordgo.Session) {
	c.RegisteredCommands = make([]*Command, 0)
	c.commands = make(map[string]Command)
	c.ResultMessages = make(chan CommandResultMessage)
	c.session = session
}

// Register associates a Command name to a Handler
func (c *Commands) Register(command Command) {
	c.RegisteredCommands = append(c.RegisteredCommands, &command)
	commandNameSplit := strings.Split(command.Name, " ")
	if len(commandNameSplit) < 1 {
		return
	}
	for _, commandName := range commandNameSplit {
		c.commands[strings.ToLower(commandName)] = command
	}
}

func (c *Commands) parse(session *discordgo.Session, commandMessage *discordgo.MessageCreate) {
	commandParsed, err := shlex.Split(commandMessage.Content, true)
	if err != nil {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Error happened " + err.Error(),
			Color:          0xb30000,
		}
		return
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Invalid Command",
			Color:          0xb30000,
		}
		return
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[strings.ToLower(commandName)]; exists {
		valid := c.Validator.Validate(&command, session, commandMessage)
		if !valid {
			c.ResultMessages <- &CommandError{
				CommandMessage: commandMessage,
				Message:        "No permissions to execute this command",
				Color:          0xb30000,
			}
			return
		}
		if commandCount < 2 {
			resultMessage := command.Handler(c.session, commandMessage, nil)
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		} else {
			resultMessage := command.Handler(c.session, commandMessage, commandParsed[1:])
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		}
	} else {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Command doesn't exist: " + commandName,
			Color:          0xb30000,
		}
	}
}

// Parse is the input sink for commands
func (c *Commands) Parse(session *discordgo.Session, commandMessage *discordgo.MessageCreate) {
	go c.parse(session, commandMessage)
}
