package main

import bolt "go.etcd.io/bbolt"

type vote2 struct {
	entity
}

func (v *vote2) Init() {
	v.entity.Init()
	v.name = "vote"
	v.onLoad = v.OnLoad
}

func (v *vote2) OnLoad(key string, val []byte, bucket *bolt.Bucket) interface{} {
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
