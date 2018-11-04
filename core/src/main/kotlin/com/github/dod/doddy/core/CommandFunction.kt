package com.github.dod.doddy.core

import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import java.lang.Exception
import kotlin.reflect.KFunction
import kotlin.reflect.KParameter
import kotlin.reflect.full.createType
import kotlin.reflect.jvm.javaType

data class CommandFunction(
        val module: Module,
        val function: KFunction<*>,
        val parameters: List<KParameter>,
        val optionals: List<Int>,
        val allArgs: Boolean,
        val commandAnnotation: Command
) {

    companion object {
        private val stringType = String::class.java
        private val intType = Int::class.java
        private val longType = Long::class.java
        private val shortType = Short::class.java
        private val doubleType = Double::class.java
    }

    fun call(event: MessageReceivedEvent, args: List<String>): CommandResult {
        if (args.size + optionals.size < parameters.size && !allArgs) {//TODO: check for too many arguments
            return InvalidArgs(args)
        }
        val params = ArrayList<Any?>()
        params.add(module)
        params.add(event)
        if (allArgs) {
            params.addAll(args)
        }
        args.forEachIndexed { index, argument ->
            val paramIndex = index + 2
            if (parameters.size > index) {
                when (parameters[index].type.javaType) {
                    stringType -> {
                        params.add(paramIndex, argument)
                    }
                    intType -> {
                        val number = argument.toIntOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    longType -> {
                        val number = argument.toLongOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    shortType -> {
                        val number = argument.toShortOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                    doubleType -> {
                        val number = argument.toDoubleOrNull()
                        if (number != null) {
                            params.add(paramIndex, number)
                        } else {
                            return InvalidArg(argument, "not a number")
                        }
                    }
                }
            }
        }
        optionals.forEach {
            if (params.size <= it || params[it] == null) {
                params.add(it, null)
            }
        }
        try {
            function.call(*params.toArray())
        } catch (exception: Exception) {
            return CommandError(exception)
        }
        return Success("bla")
    }
}