package com.github.dod.doddy

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.Feature
import com.github.dod.doddy.db.Db
import com.github.dod.doddy.db.coroutines.findOne
import com.github.dod.doddy.db.inc
import com.github.dod.doddy.guilds.Guild
import com.github.dod.doddy.users.User
import com.mongodb.Block
import com.mongodb.async.SingleResultCallback
import com.mongodb.client.model.UpdateOptions
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.eq
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlin.coroutines.suspendCoroutine

class ReputationFeature : Feature {
    @Command(names = ["reputation", "rep", "r"],
        shortDescription = "Displays a users reputation.",
        longDescription = "Returns ones own reputation unless another userName is specified.",
        docUrl = "https://github.com/Devs-On-Discord/DoDdy/tree/develop/reputation#reputation")
    suspend fun reputationCommand(event: MessageReceivedEvent, userName: String?) {
        val dbUsers = Db.instance.getCollection<User>()

        val user = parseUser(userName)   // resolves user

        val reputations = if (user == null)
            mutableMapOf()
        else
            dbUsers.findOne(User::id eq user.id)?.guildReputations ?: mutableMapOf()

        val dbGuilds = Db.instance.getCollection<Guild>()
        val dbGuildList = dbGuilds.find()

        dbGuildList.forEach({ guildIt ->
            val guildName = guildIt.name
            val localRep = reputations[guildIt.id]

            // Show reputation for each guild, put bold if the guildIt is the same guild the event channel is in
        }, { result, t -> /* TODO: Some voodoo coroutine stuff, probably similar to ReputationChanger */ })
    }
}