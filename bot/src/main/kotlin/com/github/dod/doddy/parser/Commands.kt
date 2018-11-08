package com.github.dod.doddy.parser

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.CommandFunction
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import kotlin.reflect.full.createType
import kotlin.reflect.full.findAnnotation
import kotlin.reflect.full.functions

class Commands {

    private val commandFunctions = mutableMapOf<String, CommandFunction>()

    internal val functions = mutableListOf<CommandFunction>()

    fun register(caller: Module) {
        val module = caller.javaClass.kotlin
        module.functions.forEach { function ->
            val parameters = function.parameters
            if (parameters.size > 1) {
                val commandAnnotation = function.findAnnotation<Command>()
                if (commandAnnotation != null) {
                    if (parameters[1].type == MessageReceivedEvent::class.createType()) {
                        val optionals = mutableListOf<Int>()
                        parameters.forEachIndexed { index, parameter ->
                            if (parameter.type.isMarkedNullable) {
                                optionals.add(index)
                            }
                        }
                        val functionsParams = parameters.drop(2)
                        val commandFunction = CommandFunction(
                            caller,
                            function,
                            functionsParams,
                            optionals,
                            parameters.last().type == List::class,
                            commandAnnotation
                        )
                        functions.add(commandFunction)
                        commandAnnotation.names.forEach { commandName ->
                            commandFunctions[commandName] = commandFunction
                        }
                    }
                }
            }
        }
    }

    suspend fun call(name: String, event: MessageReceivedEvent, args: List<String>): CommandResult {
        val commandFunction = commandFunctions[name] ?: return CommandNotFound(name)
        return commandFunction.call(event, args)
    }
}