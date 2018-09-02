package votes

import (
	"github.com/Devs-On-Discord/DoDdy/db"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type Vote struct {
	Id      string
	Name    string
	Message string
	Answers []Answer
}

type Answer struct {
	name    string
	emojiID string
	count   int
}

func GetVotes() ([]Vote, error) {
	votes := make([]Vote, 0)
	err := db.DB.View(func(tx *bolt.Tx) error {
		votesBucket := tx.Bucket([]byte("votes"))
		if votesBucket == nil {
			return fmt.Errorf("votes bucket doesn't exists")
		}
		votesBucket.ForEach(func(k, v []byte) error {
			voteBucket := votesBucket.Bucket(k)
			if voteBucket != nil {
				id := voteBucket.Get([]byte("id"))
				if id == nil {
					return fmt.Errorf("vote id couldn't be loaded")
				}
				name := voteBucket.Get([]byte("name"))
				if name == nil {
					return fmt.Errorf("vote name couldn't be loaded")
				}
				message := voteBucket.Get([]byte("message"))
				if id == nil {
					return fmt.Errorf("vote message couldn't be loaded")
				}
				answersBucket := voteBucket.Bucket([]byte("answers"))
				if answersBucket == nil {
					return fmt.Errorf("vote answers couldn't be loaded")
				}
				answers := make([]Answer, 0)
				answersBucket.ForEach(func(k, v []byte) error {
					answerBucket := votesBucket.Bucket(k)
					name := answerBucket.Get([]byte("name"))
					if name == nil {
						return fmt.Errorf("vote answer name couldn't be loaded")
					}
					emojiID := answerBucket.Get([]byte("emojiID"))
					if emojiID == nil {
						return fmt.Errorf("vote emojiID name couldn't be loaded")
					}
					count := answerBucket.Get([]byte("count"))
					if count == nil {
						return fmt.Errorf("vote answer count couldn't be loaded")
					}
					countInt, err := strconv.Atoi(string(count))
					if err != nil {
						return fmt.Errorf("vote answer count couldn't be converted")
					}
					answers = append(answers, Answer{name: string(name), emojiID: string(emojiID), count: countInt})
					return nil
				})
				votes = append(votes, Vote{Id: string(id), Name: string(name), Message: string(message), Answers: answers})
			}
			return nil
		})
		return nil
	})
	return votes, err
}

func AddVote(id string, name string, message string, answers []Answer) (error) {
	return db.DB.Update(func(tx *bolt.Tx) error {
		votesBucket, err := tx.CreateBucketIfNotExists([]byte("votes"))
		if err == nil {
			return fmt.Errorf("could not create votes bucket")
		}
		voteBucket, err := votesBucket.CreateBucket([]byte(id))
		if err != nil {
			return fmt.Errorf("could not create vote bucket")
		}
		err = voteBucket.Put([]byte("id"), []byte(id))
		if err != nil {
			return fmt.Errorf("could not save vote id")
		}
		err = voteBucket.Put([]byte("name"), []byte(name))
		if err != nil {
			return fmt.Errorf("could not save vote name")
		}
		err = voteBucket.Put([]byte("message"), []byte(message))
		if err != nil {
			return fmt.Errorf("could not save vote message")
		}
		answersBucket, err := voteBucket.CreateBucket([]byte("answers"))
		if err != nil {
			return fmt.Errorf("could not save vote answers")
		}
		for _, answer := range answers {
			answerBucket, err := answersBucket.CreateBucket([]byte(answer.emojiID))
			if err != nil {
				return fmt.Errorf("could not save vote answer")
			}
			err = answerBucket.Put([]byte("name"), []byte(answer.name))
			if err != nil {
				return fmt.Errorf("could not save vote answer name")
			}
			err = answerBucket.Put([]byte("emojiID"), []byte(answer.emojiID))
			if err != nil {
				return fmt.Errorf("could not save vote answer emojiID")
			}
			err = answerBucket.Put([]byte("count"), []byte(string(answer.count)))
			if err != nil {
				return fmt.Errorf("could not save vote answer count")
			}
		}
		return nil
	})
}
