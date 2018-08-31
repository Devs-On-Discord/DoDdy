package commands

type commandReply struct {
	commandResultMessage
	message string
	color   int
}
