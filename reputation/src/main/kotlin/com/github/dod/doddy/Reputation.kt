package com.github.dod.doddy

import com.github.dod.doddy.db.Db
import com.github.dod.doddy.db.coroutines.findOne
import com.github.dod.doddy.db.inc
import com.github.dod.doddy.users.User
import com.mongodb.client.model.UpdateOptions
import com.mongodb.client.result.UpdateResult
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.eq
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlin.coroutines.suspendCoroutine

// 1.0 = no bonus
private const val REPUTATION_BONUS_FACTOR_MAX = 2.0

suspend fun User.addReputation(guildId: String,
                               reputation: Long) : UpdateResult {
    val dbUsers = Db.instance.getCollection<User>()
    val dbUser = dbUsers.findOne(User::id eq this.id)

    val reputationWithBonus = dbUser
        ?.let { applyBonus(reputation, it.getGlobalReputation(), it.guildReputations[guildId] ?: 0) }
        ?: reputation

    return suspendCoroutine { continuation ->
        dbUsers.updateOne(User::id eq this.id,
            User::guildReputations.inc(guildId, reputationWithBonus),
            UpdateOptions().upsert(true)
        ) { result, throwable ->
            if (throwable != null) {
                continuation.resumeWithException(throwable)
            } else {
                continuation.resume(result)
            }
        }
    }
}

// Necessary?
suspend fun User.removeReputation(guildId: String,
                                  reputation: Long) = addReputation(guildId, -reputation)

private fun applyBonus(reputationToAdd: Long,
                       globalReputation: Long,
                       localReputation: Long): Long {
    val normalizedBonusFactor = if (globalReputation == 0L) 0.0 else 1 - (localReputation.toDouble() / globalReputation) // 0 for no bonus, 1 for complete bonus
    val bonusFactor = normalizedBonusFactor * (REPUTATION_BONUS_FACTOR_MAX - 1) + 1

    return (bonusFactor * reputationToAdd).toLong()
}