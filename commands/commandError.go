package commands

type commandError struct {
	message string
	color   int
}

func (c commandError) Message() string {
	return c.message
}

func (c commandError) Color() int {
	return c.color
}

func (c *commandError) Error() string {
	return c.message
}
