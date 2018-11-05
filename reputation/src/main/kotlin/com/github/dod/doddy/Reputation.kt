package com.github.dod.doddy

import com.github.dod.doddy.db.Db
import com.github.dod.doddy.db.incEntry
import com.github.dod.doddy.users.User
import com.mongodb.client.model.UpdateOptions
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.coroutine.findOne
import org.litote.kmongo.eq

private const val REPUTATION_BONUS_FACTOR_MAX = 2.0

suspend fun User.addReputation(guildId: String,
                               reputation: Long) {
    // TODO: Non-reliable values?
    // Computation non-random, values in database seemingly are
    val dbUsers = Db.instance.getCollection<User>()
    val dbUser = dbUsers.findOne(User::id eq this.id)

    val reputationWithBonus = if (dbUser == null) {
        reputation
    } else {
        val globalRep = dbUser.getGlobalReputation()

        val diff = globalRep - (dbUser.guildReputations[guildId] ?: 0)
        val diffPerc = diff.toDouble() / globalRep

        ((diffPerc * (REPUTATION_BONUS_FACTOR_MAX - 1) + 1) * reputation).toLong()
    }

    // Whoever this is, you're a blessing: https://stackoverflow.com/questions/25090548/push-equivalent-for-map-in-mongo
    dbUsers.updateOne(User::id eq this.id,
        incEntry(User::guildReputations, guildId, reputationWithBonus),
        UpdateOptions().upsert(true)
        ) {_, _ -> }
}

suspend fun User.removeReputation(guildId: String,
                                  reputation: Long) = addReputation(guildId, -reputation)