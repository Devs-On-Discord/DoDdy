package commands

type commandReply struct {
	message string
	color   int
}

func (c commandReply) Message() string {
	return c.message
}

func (c commandReply) Color() int {
	return c.color
}
