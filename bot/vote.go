package main

import bolt "go.etcd.io/bbolt"

type vote struct {
	entity
	title string
	message string
	answers map[string]*voteAnswer
	guild map[string]*voteGuild
}

func (v *vote) Init() {
	v.entity.Init()
	v.fields = map[string]*entityField{
		"title": {
			setter: func(val interface{}) {
				if title, ok := val.(string); ok {
					v.title = title
				}
			},
			getter: func() interface{} {
				return v.title
			},
		},
		"message": {
			setter: func(val interface{}) {
				if message, ok := val.(string); ok {
					v.message = message
				}
			},
			getter: func() interface{} {
				return v.message
			},
		},
		"answers": {
			getter: func() interface{} {
				return v.answers
			},
		},
		"guild": {
			getter: func() interface{} {
				return v.guild
			},
		},
	}
	v.name = "vote"
	v.onLoad = v.OnLoad
	v.onSave = v.OnSave
}

func (v *vote) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "title", "message":
		return string(val)
	case "answers":
		v.answers = map[string]*voteAnswer{}
		v.loadNestedBucketEntity(key, bucket, func(id string, bucket *bolt.Bucket) {
			voteAnswer := &voteAnswer{}
			voteAnswer.Init()
			voteAnswer.SetID(id)
			voteAnswer.LoadBucket(bucket)
			v.answers[id] = voteAnswer
		})
	case "guild":
		v.guild = map[string]*voteGuild{}
		v.loadNestedBucketEntity(key, bucket, func(id string, bucket *bolt.Bucket) {
			voteGuild := &voteGuild{}
			voteGuild.Init()
			voteGuild.SetID(id)
			voteGuild.LoadBucket(bucket)
			v.guild[id] = voteGuild
		})
	}
	return nil
}

func (v *vote) OnSave(key string, val interface{}, bucket *bolt.Bucket) (interface{}, error) {
	switch key {
	case "answers":
		if v.answers == nil {
			return nil, nil
		}
		v.saveNestedBucketEntities(key, bucket, len(v.answers), func(save func(entity Entity)) {
			for _, answer := range v.answers {
				save(answer)
			}
		})
	case "guild":
		println("save gui")
		if v.guild == nil {
			return nil, nil
		}
		v.saveNestedBucketEntities(key, bucket, len(v.guild), func(save func(entity Entity)) {
			for _, guild := range v.guild {
				save(guild)
			}
		})
	}
	return nil, nil
}