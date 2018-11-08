package com.github.dod.doddy.core

import net.dv8tion.jda.core.JDA

interface Feature {
    /**
     * @return the event listeners that should be added to the discord bot
     */
    fun getEventListeners(): List<Any> = emptyList()

    fun onBotReady(bot: JDA, commandFunctions: List<CommandFunction>) {}

    fun onDestroy() {}
}