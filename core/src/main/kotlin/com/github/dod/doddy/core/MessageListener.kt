package com.github.dod.doddy.core

import net.dv8tion.jda.core.events.Event
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import net.dv8tion.jda.core.hooks.EventListener
import java.util.regex.Pattern

class MessageListener(private val modules: Modules) : EventListener {
    override fun onEvent(event: Event) {
        println(event)
        when (event) {
            is MessageReceivedEvent -> onMessageReceived(event)
        }
    }

    private fun onMessageReceived(event: MessageReceivedEvent) {
        if (event.author.isBot) return
        val content = event.message.contentRaw
        val length = content.length
        if (length == 0) return
        if (content[0] == '!') {//TODO: use from db
            val parsedMessage = mutableListOf<String>()
            val m = Pattern.compile("([^\"]\\S*|\".+?\")\\s*").matcher(content.slice(1 until length))
            while (m.find()) {
                parsedMessage.add(m.group(1))
            }
            //val parsedMessage = content.slice(1 until length).trim().split("/ +/g".toRegex())
            val parsedMessageSize = parsedMessage.size
            if (parsedMessageSize == 0) return
            val commandName = parsedMessage[0]
            val commandArgs = if (parsedMessageSize > 1) { // has args
                parsedMessage.slice(1 until parsedMessageSize)
            } else {
                emptyList()
            }

            val messageAction = when (
                val commandResult = modules.onCommand(commandName, event, commandArgs)) {
                is CommandNotFound -> event.channel.sendMessage("command not found :" + commandResult.commandName)
                is InvalidArgs -> event.channel.sendMessage(commandResult.args.toString())
                is InvalidArg -> event.channel.sendMessage(commandResult.arg + " " + commandResult.error)
                is CommandError -> {
                    commandResult.exception.printStackTrace()
                    event.channel.sendMessage(commandResult.exception.toString())
                }
                else -> null
            }
            messageAction?.queue()
        }
    }
}