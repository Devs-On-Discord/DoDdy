package com.github.dod.doddy.core

import kotlin.reflect.KClass
import kotlin.reflect.KFunction
import kotlin.reflect.KParameter

data class CommandFunction(
    val module: KClass<out Module>,
    val function: KFunction<*>,
    val parameters: List<KParameter>
) {
    fun call(args: List<String>): CommandResult {
        if (args.size != parameters.size) {
            return InvalidArgs(args)
        }
        val params = ArrayList<Any>(args.size)
        args.forEachIndexed { index, argument ->
            when (parameters[index].type) {
                String::class -> {
                    params[index] = argument
                }
                Int::class -> {
                    val number = argument.toIntOrNull()
                    if (number != null) {
                        params[index] = number
                    } else {
                        return InvalidArg(argument, "not a number")
                    }
                }
                Long::class -> {
                    val number = argument.toLongOrNull()
                    if (number != null) {
                        params[index] = number
                    } else {
                        return InvalidArg(argument, "not a number")
                    }
                }
                Short::class -> {
                    val number = argument.toShortOrNull()
                    if (number != null) {
                        params[index] = number
                    } else {
                        return InvalidArg(argument, "not a number")
                    }
                }
                Double::class -> {
                    val number = argument.toDoubleOrNull()
                    if (number != null) {
                        params[index] = number
                    } else {
                        return InvalidArg(argument, "not a number")
                    }
                }
            }
        }
        function.call(params)
        return Success("bla")
    }
}