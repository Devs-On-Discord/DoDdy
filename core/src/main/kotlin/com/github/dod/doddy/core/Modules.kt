package com.github.dod.doddy.core

import java.util.*

class Modules internal constructor() {
    private val commands = mutableMapOf<String, Command>()

    internal val modules = LinkedList<Module>()

    fun add(module: Module) {
        modules.add(module)
        module.getCommands().forEach { command ->
            command.names.forEach { name ->
                commands[name] = command
            }
        }
    }

    internal fun onCommand(name: String, args: List<String>): Boolean {
        return commands[name]?.onCommand(args) ?: false
    }
}