package com.github.dod.doddy.parser

import com.github.dod.doddy.core.Feature
import net.dv8tion.jda.core.JDA
import net.dv8tion.jda.core.JDABuilder
import java.io.File

object Discord {
    private val modules = Features()

    val bot: JDA = JDABuilder(File("../discord.token").readText()).build()

    init {
        bot.addEventListener(MessageListener(modules))
    }

    fun add(feature: Feature) {
        modules.add(feature)
        feature.getEventListeners().forEach { eventListener ->
            bot.addEventListener(eventListener)
        }
    }

    fun start() {
        modules.commandsReady()
        bot.awaitReady()
    }

    fun stop() {
        modules.destroy()
        bot.shutdown()
    }
}