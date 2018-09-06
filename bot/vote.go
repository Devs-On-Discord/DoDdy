package main

import bolt "go.etcd.io/bbolt"

type vote struct {
	entity
}

func (v *vote) Init() {
	v.entity.Init()
	v.name = "vote"
	v.onLoad = v.OnLoad
}

func (v *vote) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
	switch key {
	case "name", "message":
		return string(val)
	case "answers":
		answers := map[string]Entity{}
		if answersBucket := bucket.Bucket([]byte(key)); answersBucket != nil {
			answersCursor := answersBucket.Cursor()
			for k, _ := answersCursor.First(); k != nil; k, _ = answersCursor.Next() {
				if answerBucket := answersBucket.Bucket(k); answerBucket != nil {
					answer := &voteAnswer{}
					answer.Init()
					answer.id = string(k)
					answer.LoadBucket(answerBucket)
					answers[answer.id] = answer
				}
			}
		}
		return answers
	case "guild":
		guilds := map[string]Entity{}
		if guildsBucket := bucket.Bucket([]byte(key)); guildsBucket != nil {
			guildsCursor := guildsBucket.Cursor()
			for k, _ := guildsCursor.First(); k != nil; k, _ = guildsCursor.Next() {
				if guildBucket := guildsBucket.Bucket(k); guildBucket != nil {
					guild := &voteGuild{}
					guild.Init()
					guild.id = string(k)
					guild.LoadBucket(guildBucket)
					guilds[guild.id] = guild
				}
			}
		}
		return guilds
	}
	return nil
}
