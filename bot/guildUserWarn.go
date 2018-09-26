package main

type guildUserWarn struct {
	entity
	reason string
	authorID string
	timestamp uint64
}