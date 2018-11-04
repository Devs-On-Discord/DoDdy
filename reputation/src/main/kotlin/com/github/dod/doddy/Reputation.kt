package com.github.dod.doddy

import com.github.dod.doddy.db.Db
import com.github.dod.doddy.db.incEntry
import com.github.dod.doddy.users.User
import com.mongodb.client.model.UpdateOptions
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.coroutine.findOne
import org.litote.kmongo.eq

suspend fun User.addReputation(guildId: String,
                       reputation: Long) {
    val dbUsers = Db.instance.getCollection<User>()
    // Next line makes this a coroutine, otherwise this method doesn't work (Don't know why)
    val dbUser = dbUsers.findOne()

    // Whoever this is, you're a blessing: https://stackoverflow.com/questions/25090548/push-equivalent-for-map-in-mongo
    dbUsers.updateOne(User::id eq this.id,
        incEntry(User::guildReputations, guildId, reputation),
        UpdateOptions().upsert(true)
        ) {_, _ -> }
}

suspend fun User.removeReputation(guildId: String,
                          reputation: Long) = addReputation(guildId, -reputation)