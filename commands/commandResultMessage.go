package commands

type commandResultMessage interface {
	Message() (string)
	Color() (int)
}
