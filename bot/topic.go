package main

import (
	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

type topic struct {
	entity
	title     string
	color     string // Format 0x{HexColor} without #
	iconURL   string
	topicID   string
	channelID string
}

func (t *topic) Init() {
	t.entity.Init()
	t.fields = map[string]*entityField{
		"title": {
			setter: func(val interface{}) {
				if title, ok := val.(string); ok {
					t.title = title
				}
			},
			getter: func() interface{} {
				return t.title
			},
		},
		"color": {
			setter: func(val interface{}) {
				if color, ok := val.(string); ok {
					t.color = color
				}
			},
			getter: func() interface{} {
				return t.color
			},
		},
		"iconURL": {
			setter: func(val interface{}) {
				if iconURL, ok := val.(string); ok {
					t.iconURL = iconURL
				}
			},
			getter: func() interface{} {
				return t.iconURL
			},
		},
		"topicID": {
			setter: func(val interface{}) {
				if topicID, ok := val.(string); ok {
					t.topicID = topicID
				}
			},
			getter: func() interface{} {
				return t.topicID
			},
		},
		"channelID": {
			setter: func(val interface{}) {
				if channelID, ok := val.(string); ok {
					t.channelID = channelID
				}
			},
			getter: func() interface{} {
				return t.channelID
			},
		},
	}
	t.name = "topic"
	t.onLoad = t.OnLoad
}

func (t *topic) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "title", "color", "icon", "topicID", "channelID":
		return string(val)
	}
	return nil
}

func (t *topic) CreateQuestion(session *discordgo.Session, message discordgo.MessageCreate) {
	content, err := message.ContentWithMoreMentionsReplaced(session)
	if err != nil {
		return
	}
	message.Content = content

}
