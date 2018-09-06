package main

import (
	"bytes"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var (
	guilds                      = []byte("guilds")
	guildName                   = []byte("name")
	guildPrefix                 = []byte("prefix")
	guildAnnouncementsChannelID = []byte("announcementsChannelID")
	guildVotesChannelID         = []byte("votesChannelID")
	guildRoles                  = []byte("roles")
)

// Guilds is an object able to access the database, and possesses a list of guilds
type Guilds struct {
	db     *bolt.DB
	Guilds map[string]*Guild
}

// Guild contains cached data from the database, as well as an access to said database
type Guild struct {
	db                     *bolt.DB
	id                     string
	name                   string
	Prefix                 string
	AnnouncementsChannelID string
	VotesChannelID         string
	Roles                  map[Role]string // Value: guild specific role id
}

const (
	notSetup         = "bot isn't set up for this guild"
	bucketNotCreated = "guild's bucket couldn't be created"
)

// Init constructs the Guilds object
func (g *Guilds) Init(db *bolt.DB) {
	g.db = db
	g.Guilds = make(map[string]*Guild)
}

// Create adds a guild to the database, and associates it to a name
func (g *Guilds) Create(id string, name string) error {
	err := g.db.Update(func(tx *bolt.Tx) error {
		guildsBucket, err := tx.CreateBucketIfNotExists([]byte(guilds))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		guildBucket, err := guildsBucket.CreateBucket([]byte(id))
		if err != nil {
			return fmt.Errorf(bucketNotCreated)
		}
		err = guildBucket.Put(guildName, []byte(name))
		if err != nil {
			return fmt.Errorf("name couldn't be saved")
		}
		err = guildBucket.Put(guildPrefix, []byte("!"))
		if err != nil {
			return fmt.Errorf("prefix couldn't be saved")
		}
		return nil
	})
	if err == nil {
		guild := &Guild{db: g.db, id: id, name: name, Roles: make(map[Role]string)}
		g.Guilds[id] = guild
	}
	return err
}

func (g *Guilds) loadGuild(guildsBucket *bolt.Bucket, guildID string) *Guild {
	guildBucket := guildsBucket.Bucket([]byte(guildID))
	if guildBucket != nil {
		guildCursor := guildBucket.Cursor()
		guild := &Guild{db: g.db, id: string(guildID), Prefix: "", Roles: make(map[Role]string)}
		for k, v := guildCursor.First(); k != nil; k, v = guildCursor.Next() {
			if bytes.Equal(k, guildName) {
				guild.name = string(v)
			} else if bytes.Equal(k, guildPrefix) {
				guild.Prefix = string(v)
			} else if bytes.Equal(k, guildAnnouncementsChannelID) {
				guild.AnnouncementsChannelID = string(v)
			} else if bytes.Equal(k, guildVotesChannelID) {
				guild.VotesChannelID = string(v)
			} else if bytes.Equal(k, guildRoles) {
				rolesBucket := guildBucket.Bucket(guildRoles)
				if rolesBucket != nil {
					rolesCursor := rolesBucket.Cursor()
					for k, v := rolesCursor.First(); k != nil; k, v = rolesCursor.Next() {
						roleInt, err := strconv.Atoi(string(k))
						if err == nil {
							if role, exists := RoleInt[roleInt]; exists {
								guild.Roles[role] = string(v)
							}
						}
					}
				}
			}
		}
		g.Guilds[guild.id] = guild
		return guild
	}
	return nil
}

// LoadGuilds caches the guilds from the database
func (g *Guilds) LoadGuilds() *Guilds {
	g.db.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket(guilds)
		guildsCursor := guildsBucket.Cursor()
		for k, _ := guildsCursor.First(); k != nil; k, _ = guildsCursor.Next() {
			if _, exists := g.Guilds[string(k)]; !exists {
				g.loadGuild(guildsBucket, string(k))
			}
		}
		return nil
	})
	return g
}

// Guild returns a cached guild instance, or pulls it from the database
func (g *Guilds) Guild(id string) (*Guild, error) {
	guild, exists := g.Guilds[id]
	if exists {
		return guild, nil
	}
	err := g.db.View(func(tx *bolt.Tx) error {
		guildsBucket := tx.Bucket(guilds)
		if guildsBucket != nil {
			guild = g.loadGuild(guildsBucket, id)
		}
		return nil
	})
	return guild, err

}

func (g *Guild) bucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	guildsBucket := tx.Bucket(guilds)
	if guildsBucket == nil {
		println("guildsBucket==nil")
		return nil, fmt.Errorf(bucketNotCreated)
	}
	guildBucket := guildsBucket.Bucket([]byte(g.id))
	if guildBucket == nil {
		return nil, fmt.Errorf(notSetup)
	}
	return guildBucket, nil
}

func (g *Guild) set(key []byte, value []byte) error {
	err := g.db.Update(func(tx *bolt.Tx) error {
		bucket, err := g.bucket(tx)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return fmt.Errorf("%s couldn't be saved", string(key))
		}
		return nil
	})
	return err
}

// SetPrefix defines the prefix of a single guild
func (g *Guild) SetPrefix(prefix string) error {
	err := g.set(guildPrefix, []byte(prefix))
	if err == nil {
		g.Prefix = prefix
	}
	return err
}

// SetAnnouncementsChannel defines the announcement channel for a single guild
func (g *Guild) SetAnnouncementsChannel(channelID string) error {
	err := g.set(guildAnnouncementsChannelID, []byte(channelID))
	if err == nil {
		g.AnnouncementsChannelID = channelID
	}
	return err
}

// SetVotesChannel defines the voting channel for a single guild
func (g *Guild) SetVotesChannel(channelID string) error {
	err := g.set(guildVotesChannelID, []byte(channelID))
	if err == nil {
		g.VotesChannelID = channelID
	}
	return err
}

// SetRole assigns a role on the server to one of the permission levels
func (g *Guild) SetRole(name string, id string) error {
	if role, exists := CommandRoleNames[name]; exists {
		err := g.db.Update(func(tx *bolt.Tx) error {
			bucket, err := g.bucket(tx)
			if err != nil {
				return err
			}
			if rolesBucket, err := bucket.CreateBucketIfNotExists(guildRoles); err == nil {
				err = rolesBucket.Put([]byte(strconv.Itoa(int(role))), []byte(id))
				if err != nil {
					return err
				}
			} else {
				return err
			}
			return nil
		})
		if err == nil {
			g.Roles[role] = id
		}
		return err
	}
	return fmt.Errorf("role %s doesn't exists\nEnter !roles to see possible names", string(name))

}
