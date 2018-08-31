package commands

type Handle func(args []string) (commandResultMessage)

type Command struct {
	Name string
	Handler Handle
}
