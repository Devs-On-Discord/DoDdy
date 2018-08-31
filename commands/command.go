package commands

type Handle func(args []string)

type Command struct {
	Name string
	Handler Handle
	channel chan <- []string
}
