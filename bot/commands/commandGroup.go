package commands

type CommandGroup interface {
	Commands() []*Command
}
