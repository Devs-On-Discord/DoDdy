package main

import bolt "go.etcd.io/bbolt"

type topic struct {
	entity
	title   string
	color   string // Format 0x{HexColor} without #
	iconURL    string
	topicID string
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
	}
	t.name = "topic"
	t.onLoad = t.OnLoad
}

func (t *topic) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "title", "color", "icon", "topicID":
		return string(val)
	}
	return nil
}
