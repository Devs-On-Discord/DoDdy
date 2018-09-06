package main

import (
	"bytes"
	"strconv"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

var (
	votes              = []byte("votes")
	voteName           = []byte("name")
	voteMessage        = []byte("message")
	voteAnswers        = []byte("answers")
	voteGuilds         = []byte("guilds")
	voteGuildChannelID = []byte("channelID")
	voteGuildMessageID = []byte("messageID")
	voteAnswerName     = []byte("name")
	voteAnswerCount    = []byte("count")
)

//TODO: override message with messageID and channelID with editing when same id is used
//TODO: add single server votes
//TODO: add initial reactions for each answer to an vote message

// Votes is a local cache to the votes in the database
type Votes struct {
	db      *bolt.DB
	session *discordgo.Session
	Votes   map[string]*Vote
	// Is there for faster calculations in reaction add, remove
	channelVotes map[string]*Vote // Key: channelID
}

// Vote references a single vote in the cache
type Vote struct {
	db      *bolt.DB
	ID      string
	Name    string
	Message string
	Answers map[string]*Answer    // Key: emojiID
	Guilds  map[string]*GuildVote // Key: guildID
}

// GuildVote defines where a vote is located
type GuildVote struct {
	MessageID string
	ChannelID string
}

// Answer contains the number of answers given with a certain emoji
type Answer struct {
	name    string
	emojiID string
	count   int
}

// Init constructs the Votes object
func (v *Votes) Init(db *bolt.DB, session *discordgo.Session) {
	v.db = db
	v.session = session
	v.session.AddHandler(v.reactionAdded)
	v.session.AddHandler(v.reactionRemoved)
	v.Votes = make(map[string]*Vote)
	v.channelVotes = make(map[string]*Vote)
	// For faster calculation in reaction add, remove we need all votes
	v.LoadVotes()
	v.fillChannelVotes()
}

// Create a vote
func (v *Votes) Create(id string, name string, message string, answers map[string]*Answer, guilds map[string]*GuildVote) error {
	err := v.db.Update(func(tx *bolt.Tx) error {
		votesBucket, err := tx.CreateBucketIfNotExists(votes)
		if err != nil {
			return err
		}
		voteBucket, err := votesBucket.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return err
		}
		err = voteBucket.Put(voteName, []byte(name))
		if err != nil {
			return err
		}
		err = voteBucket.Put(voteMessage, []byte(message))
		if err != nil {
			return err
		}
		voteAnswersBucket, err := voteBucket.CreateBucketIfNotExists(voteAnswers)
		if err != nil {
			return err
		}
		voteGuildsBucket, err := voteBucket.CreateBucketIfNotExists(voteGuilds)
		if err != nil {
			return err
		}
		for _, answer := range answers {
			voteAnswerBucket, err := voteAnswersBucket.CreateBucketIfNotExists([]byte(answer.emojiID))
			if err != nil {
				return err
			}
			err = voteAnswerBucket.Put(voteAnswerName, []byte(answer.name))
			if err != nil {
				return err
			}
			err = voteAnswerBucket.Put(voteAnswerCount, []byte(string(answer.count)))
			if err != nil {
				return err
			}
		}
		for guildID, guild := range guilds {
			voteGuildBucket, err := voteGuildsBucket.CreateBucketIfNotExists([]byte(guildID))
			if err != nil {
				return err
			}
			err = voteGuildBucket.Put(voteGuildMessageID, []byte(guild.MessageID))
			if err != nil {
				return err
			}
			err = voteGuildBucket.Put(voteGuildChannelID, []byte(guild.ChannelID))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		vote := &Vote{ID: id, Name: name, Message: message, Answers: answers, Guilds: guilds}
		v.Votes[id] = vote
		v.fillChannelVotesForVote(vote)
	}
	return err
}

func (v *Votes) loadVote(votesBucket *bolt.Bucket, id string) *Vote {
	voteBucket := votesBucket.Bucket([]byte(id))
	if voteBucket != nil {
		voteCursor := voteBucket.Cursor()
		vote := &Vote{db: v.db, ID: string(id), Answers: make(map[string]*Answer), Guilds: make(map[string]*GuildVote)}
		for k, v := voteCursor.First(); k != nil; k, v = voteCursor.Next() {
			if bytes.Equal(k, voteName) {
				vote.Name = string(v)
			} else if bytes.Equal(k, voteMessage) {
				vote.Message = string(v)
			} else if bytes.Equal(k, voteAnswers) {
				answersBucket := voteBucket.Bucket(k)
				answersCursor := answersBucket.Cursor()
				for k, _ := answersCursor.First(); k != nil; k, _ = answersCursor.Next() {
					answer := &Answer{emojiID: string(k)}
					answerBucket := answersBucket.Bucket(k)
					if answersBucket == nil {
						continue
					}
					answerCursor := answerBucket.Cursor()
					for k, v := answerCursor.First(); k != nil; k, v = answerCursor.Next() {
						if bytes.Equal(k, voteAnswerName) {
							answer.name = string(v)
						} else if bytes.Equal(k, voteAnswerCount) {
							countInt, err := strconv.Atoi(string(v))
							if err == nil {
								answer.count = countInt
							} else {
								answer.count = 0
							}
						}
					}
					vote.Answers[answer.emojiID] = answer
				}
			} else if bytes.Equal(k, voteGuilds) {
				guildsBucket := voteBucket.Bucket(k)
				guildsCursor := guildsBucket.Cursor()
				for k, _ := guildsCursor.First(); k != nil; k, _ = guildsCursor.Next() {
					guildBucket := guildsBucket.Bucket(k)
					if guildBucket == nil {
						continue
					}
					guildCursor := guildBucket.Cursor()
					guildVote := &GuildVote{}
					for k, v := guildCursor.First(); k != nil; k, v = guildCursor.Next() {
						if bytes.Equal(k, voteGuildChannelID) {
							guildVote.ChannelID = string(v)
						} else if bytes.Equal(k, voteGuildMessageID) {
							guildVote.MessageID = string(v)
						}
					}
					vote.Guilds[string(k)] = guildVote
				}
			}
		}
		v.Votes[id] = vote
		return vote
	}
	return nil
}

// Vote fetches and returns a vote by ID
func (v *Votes) Vote(id string) (*Vote, error) {
	vote, exists := v.Votes[id]
	if exists {
		return vote, nil
	}
	err := v.db.View(func(tx *bolt.Tx) error {
		votesBucket := tx.Bucket(votes)
		if votesBucket != nil {
			vote = v.loadVote(votesBucket, id)
		}
		return nil
	})
	return vote, err

}

// LoadVotes fetches the votes in the database to the cache
func (v *Votes) LoadVotes() *Votes {
	v.db.View(func(tx *bolt.Tx) error {
		votesBucket := tx.Bucket(votes)
		votesCursor := votesBucket.Cursor()
		for k, _ := votesCursor.First(); k != nil; k, _ = votesCursor.Next() {
			if _, exists := v.Votes[string(k)]; !exists {
				v.loadVote(votesBucket, string(k))
			}
		}
		return nil
	})
	return v
}

func (v *Votes) fillChannelVotesForVote(vote *Vote) {
	for _, guild := range vote.Guilds {
		v.channelVotes[guild.ChannelID] = vote
	}
}

func (v *Votes) fillChannelVotes() {
	for _, vote := range v.Votes {
		v.fillChannelVotesForVote(vote)
	}
}

func (v *Votes) reactionAdded(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if vote, exists := v.channelVotes[reaction.ChannelID]; exists {
		for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				//IncreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}
		go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}

func (v *Votes) reactionRemoved(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if session.State.User.ID == reaction.UserID { // Ignore bot reactions
		return
	}
	if vote, exists := v.Votes[reaction.ChannelID]; exists {
		for _, answer := range vote.Answers {
			if answer.emojiID == reaction.Emoji.ID {
				//DecreaseVoteAnswer(vote.Id, answer.emojiID)
				break
			}
		}
		go session.MessageReactionsRemoveAll(reaction.ChannelID, reaction.MessageID)
	}
}
