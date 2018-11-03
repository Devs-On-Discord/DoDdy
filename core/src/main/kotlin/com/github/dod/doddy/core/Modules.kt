package com.github.dod.doddy.core

import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import java.util.*

class Modules internal constructor() {
    private val commands = Commands()

    private val modules = LinkedList<Module>()

    fun add(module: Module) {
        modules.add(module)
        commands.register(module)
    }

    internal fun onCommand(name: String, event: MessageReceivedEvent, args: List<String>): CommandResult {
        return commands.call(name, event, args)
    }
}