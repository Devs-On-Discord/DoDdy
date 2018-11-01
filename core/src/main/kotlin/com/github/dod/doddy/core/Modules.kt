package com.github.dod.doddy.core

object Modules {
    private var modules = mutableListOf<Module>()

    private var commandModules = mutableMapOf<String, Module>()

    fun add(module: Module) {
        modules.add(module)
        module.getCommands().forEach {
            commandModules[it] = module
        }
    }

    internal fun onCommand(command: String, args: Array<Any>) {
        val module = commandModules[command]
        if (module != null) {
            module.onCommand(command, args)
        } else {
            modules.forEach { it.onCommand(command, args) }
        }
    }
}