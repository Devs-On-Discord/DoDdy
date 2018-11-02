package com.github.dod.doddy.core

object Modules {
    private val commands = mutableMapOf<String, Command>()

    fun add(module: Module) {
        module.getCommands().forEach { command ->
            command.names.forEach {name ->
                commands[name] = command
            }
        }
    }

    internal fun onCommand(name: String, args: Array<Any>): Boolean {
        return commands[name]?.onCommand(args) ?: false
    }
}