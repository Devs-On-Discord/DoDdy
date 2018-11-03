package com.github.dod.doddy.core

import net.dv8tion.jda.core.JDA
import net.dv8tion.jda.core.JDABuilder
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import net.dv8tion.jda.core.hooks.ListenerAdapter


object Discord : ListenerAdapter() {
    private val modules = Modules()

    val bot: JDA = JDABuilder("token").build()

    init {

    }

    fun add(module: Module) {
        modules.add(module)
        module.getEventListeners().forEach { eventListener ->
            bot.addEventListener(eventListener)
        }
    }

    fun start() {
        bot.awaitReady()
    }

    override fun onMessageReceived(event: MessageReceivedEvent) {
        if (event.author.isBot) return
        val content = event.message.contentRaw
        val length = content.length
        if (length == 0) return
        if (content[0] == '!') {//TODO: use from db
            val parsedMessage = content.slice(1 until length - 1).trim().split("/ +/g")
            val parsedMessageSize = parsedMessage.size
            if (parsedMessageSize == 0) return
            val commandName = parsedMessage[0]
            val commandArgs = if (parsedMessageSize > 1) { // has args
                parsedMessage.subList(1, parsedMessageSize - 1)
            } else {
                emptyList()
            }
            modules.onCommand(commandName, commandArgs)
        }
    }
}