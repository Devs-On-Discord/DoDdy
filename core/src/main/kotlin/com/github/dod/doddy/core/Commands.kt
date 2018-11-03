package com.github.dod.doddy.core

import net.dv8tion.jda.core.entities.User
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import kotlin.reflect.KClass
import kotlin.reflect.full.memberFunctions

class Commands {

    private val commandFunctions = mutableMapOf<String, CommandFunction>()

    fun register(module: KClass<out Module>) {
        module.memberFunctions.forEach { function ->
            val parameters = function.parameters
            if (parameters.isEmpty()) {
                val commandAnnotation = function.annotations.find { annotation -> annotation is Command }
                if (commandAnnotation != null && commandAnnotation is Command) {
                    commandAnnotation.names.forEach { commandName ->
                        if (parameters.first().type == MessageReceivedEvent::class) {
                            commandFunctions[commandName] = CommandFunction(
                                module,
                                function,
                                parameters.drop(0),
                                parameters.last().type == List::class
                            )
                        }
                    }
                }
            }
        }
    }

    fun call(name: String, event: MessageReceivedEvent, args: List<String>): CommandResult {
        val commandFunction = commandFunctions[name] ?: return CommandNotFound(name)
        return commandFunction.call(event, args)
    }
}