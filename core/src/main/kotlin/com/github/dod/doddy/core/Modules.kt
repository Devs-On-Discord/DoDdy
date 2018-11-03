package com.github.dod.doddy.core

import java.util.*

class Modules internal constructor() {
    private val commands = Commands()

    private val modules = LinkedList<Module>()

    fun add(module: Module) {
        modules.add(module)
        commands.register(module.javaClass.kotlin)
    }

    internal fun onCommand(name: String, args: List<String>): CommandResult {
        return commands.call(name, args)
    }
}