package com.github.dod.doddy.core

import net.dv8tion.jda.core.JDA
import net.dv8tion.jda.core.JDABuilder

object Discord {
    private val modules = Modules()

    val bot: JDA = JDABuilder("token").build()

    init {
        bot.addEventListener(MessageListener(modules))
    }

    fun add(module: Module) {
        modules.add(module)
        module.getEventListeners().forEach { eventListener ->
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