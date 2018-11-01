package com.github.dod.doddy.core

object Modules {
    private val moduleCommands = mutableMapOf<Command, Module>()
    private val commandNames = mutableMapOf<String, Command>()

    fun add(module: Module) {
        module.getCommands().forEach { command ->
            moduleCommands[command] = module
            command.names.forEach {name ->
                commandNames[name] = command
            }
        }
    }

    internal fun onCommand(name: String, args: Array<Any>): Boolean {
        return commandNames[name]?.onCommand(args) ?: false
    }
}