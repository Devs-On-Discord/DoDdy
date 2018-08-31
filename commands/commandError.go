package commands

type commandError struct {
	commandResultMessage
	message string
	color   int
}

func (c *commandError) Error() string {
	return c.message
}
