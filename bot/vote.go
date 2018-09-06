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
	println("onLoad vote", key)
	switch key {
	case "name", "message":
		return string(val)
	case "answers":
		var answers []*voteAnswer
		if answersBucket := bucket.Bucket([]byte(key)); answersBucket != nil {
			answersCursor := answersBucket.Cursor()
			for k, _ := answersCursor.First(); k != nil; k, _ = answersCursor.Next() {
				answer := &voteAnswer{}
				answer.id = string(k)
				answer.LoadBucket(answersBucket)
				answers = append(answers, answer)
			}
		}
		return answers
	}
	return nil
}
