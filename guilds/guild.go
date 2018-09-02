package guilds

import (
	bolt "go.etcd.io/bbolt"
	"fmt"
	"github.com/Devs-On-Discord/DoDdy/db"
)

type guild struct {
	id string
	name                   string
	announcementsChannelID string
}

func GetAnnouncementChannels() ([]string, error) {
	channels := make([]string, 0)
	err := db.Db.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket([]byte("guilds"))
		if guildsBucket == nil {
			return fmt.Errorf("bot isn't setted up for this guild")
		}
		guildsBucket.ForEach(func(k, v []byte) error {
			println("key " + string(k))
			guildBucket := guildsBucket.Bucket(k)
			if guildBucket != nil {
				announcementsChannelID := guildBucket.Get([]byte("announcementsChannelID"))
				println("id " + string(announcementsChannelID))
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

func SetAnnouncementsChannel(guildId string, channelId string) (error) {
	return db.Db.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			return fmt.Errorf("guilds bucket couldn't be created")
		}
		guildBucket := guildsBucket.Bucket([]byte(guildId))
		if guildBucket == nil {
			return fmt.Errorf("bot isn't setted up for this guild")
		}
		err = guildBucket.Put([]byte("announcementsChannelID"), []byte(channelId))
		if err != nil {
			return fmt.Errorf("announcements channel id couldn't be saved")
		}
		return nil
	})
}

func Create(guildId string, name string) (error) {
	return db.Db.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			return fmt.Errorf("guilds bucket couldn't be created")
		}
		guildBucket, err := guildsBucket.CreateBucket([]byte(guildId))
		if err != nil {
			return fmt.Errorf("guild bucket couldn't be created")
		}
		err = guildBucket.Put([]byte("name"), []byte(name))
		if err != nil {
			return fmt.Errorf("name couldn't be saved")
		}
		return nil
	})
}
