package commands

import "fmt"

type commandError struct {
	message string
}

func (c *commandError) Error() string {
	return fmt.Sprintf("%s", c.message)
}
