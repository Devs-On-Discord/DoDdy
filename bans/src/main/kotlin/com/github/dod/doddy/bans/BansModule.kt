package com.github.dod.doddy.bans

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.entities.Member
import net.dv8tion.jda.core.events.message.MessageReceivedEvent

class BansModule : Module {

    @Command(names = ["ban"],
        shortDescription = "Bans the specified member.",
        longDescription = "Bans the specified member and provides an option for other node mods " +
            "to also ban the user in their nodes.",
        docUrl = "https://github.com/Devs-On-Discord/DoDdy/blob/develop/bans/README.md#bans")
    fun ban(event: MessageReceivedEvent, member: Member) {
        val test = EmbedBuilder()
            .setColor(0xFFFFFF)
            .setTitle("Success")
            .setDescription("You successfully banned ${member.user.name}")
            .build()
        event.channel.sendMessage(test).queue()
    }
}