package main

/*func deletePrefix(guildID string) error {
	if _, ok := prefixes[guildID]; ok {
		delete(prefixes, guildID)
	}
	return db.Update(func(tx *bolt.Tx) error {
		nodeBucket, err := tx.CreateBucketIfNotExists([]byte("Nodes"))
		if err != nil {
			return err
		}
		guildBucket, err := nodeBucket.CreateBucketIfNotExists([]byte(guildID))
		if err != nil {
			return err
		}
		if guildBucket.Delete([]byte("Prefix")) != nil {
			return err
		}
		return nil
	})
}

func setPrefix(guildID, prefix string) error {
	prefixes[guildID] = prefix
	return db.Update(func(tx *bolt.Tx) error {
		nodeBucket, err := tx.CreateBucketIfNotExists([]byte("Nodes"))
		if err != nil {
			return err
		}
		guildBucket, err := nodeBucket.CreateBucketIfNotExists([]byte(guildID))
		if err != nil {
			return err
		}
		if guildBucket.Put([]byte("Prefix"), []byte(prefix)) != nil {
			return err
		}
		return nil
	})
}*/
// Rewrite this code to use non-nested nodes since apparently we're doing that now
