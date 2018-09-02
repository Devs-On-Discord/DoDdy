package guilds

import (
	"fmt"

	"github.com/Devs-On-Discord/DoDdy/db"
	bolt "go.etcd.io/bbolt"
)

type guild struct {
	id                     string
	name                   string
	announcementsChannelID string
	votesChannelID         string
}

const (
	notSetup         = "bot isn't set up for this guild"
	bucketNotCreated = "guild's bucket couldn't be created"
)

// GetAnnouncementChannels returns the anouncement channels for every guild
func GetAnnouncementChannels() ([]string, error) {
	channels := make([]string, 0)
	err := db.DB.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket([]byte("guilds"))
		if guildsBucket == nil {
			return fmt.Errorf(notSetup)
		}
		guildsBucket.ForEach(func(k, v []byte) error {
			guildBucket := guildsBucket.Bucket(k)
			if guildBucket != nil {
				announcementsChannelID := guildBucket.Get([]byte("announcementsChannelID"))
				if announcementsChannelID != nil {
					channels = append(channels, string(announcementsChannelID))
				}
			}
			return nil
		})
		return nil
	})
	return channels, err
}

// SetAnnouncementsChannel sets the anouncement channels for a single guild
func SetAnnouncementsChannel(guildID string, channelID string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		guildBucket := guildsBucket.Bucket([]byte(guildID))
		if guildBucket == nil {
			return fmt.Errorf(notSetup)
		}
		err = guildBucket.Put([]byte("announcementsChannelID"), []byte(channelID))
		if err != nil {
			return fmt.Errorf("announcements channel ID couldn't be saved")
		}
		return nil
	})
}

// GetVotesChannels returns the vote channels for every guild
func GetVotesChannels() ([]string, error) {
	channels := make([]string, 0)
	err := db.DB.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket([]byte("guilds"))
		if guildsBucket == nil {
			return fmt.Errorf(notSetup)
		}
		guildsBucket.ForEach(func(k, v []byte) error {
			guildBucket := guildsBucket.Bucket(k)
			if guildBucket != nil {
				votesChannelID := guildBucket.Get([]byte("votesChannelID"))
				if votesChannelID != nil {
					channels = append(channels, string(votesChannelID))
				}
			}
			return nil
		})
		return nil
	})
	return channels, err
}

// SetVotesChannel sets the vote channel for a specific guild
func SetVotesChannel(guildID string, channelID string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		guildBucket := guildsBucket.Bucket([]byte(guildID))
		if guildBucket == nil {
			return fmt.Errorf(notSetup)
		}
		err = guildBucket.Put([]byte("votesChannelID"), []byte(channelID))
		if err != nil {
			return fmt.Errorf("votes channel ID couldn't be saved")
		}
		return nil
	})
}

// Create adds a guild to the database
func Create(guildID string, name string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		guildBucket, err := guildsBucket.CreateBucket([]byte(guildID))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		err = guildBucket.Put([]byte("name"), []byte(name))
		if err != nil {
			return fmt.Errorf("name couldn't be saved")
		}
		return nil
	})
}
