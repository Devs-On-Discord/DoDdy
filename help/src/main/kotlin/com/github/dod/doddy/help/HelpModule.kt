package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.events.message.MessageReceivedEvent

class HelpModule : Module {
    @Command(names = ["help", "h"])
    fun help(event: MessageReceivedEvent) {
        event.channel.sendMessage("@everyone").queue()
    }
}