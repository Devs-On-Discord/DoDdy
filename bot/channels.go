package main

type Channel int

var (
	CommandChannelNames = map[string]Channel{
		"announcements": Announcements,
		"votes":         Votes,
	}
	ChannelInt = map[int]Channel{
		0: Announcements,
		1: Votes,
	}
)

const (
	Announcements Channel = 0
	Votes         Channel = 1
)
