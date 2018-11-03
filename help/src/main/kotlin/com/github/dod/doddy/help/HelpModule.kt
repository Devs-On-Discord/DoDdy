package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.entities.User

class HelpModule : Module {
    override fun getCommands() = setOf(
            HelpCommand())

    override fun getEventListeners(): List<Any> {
        return emptyList()
    }

    @Command(names = ["help", "h"])
    fun help(user: User) {

    }
}