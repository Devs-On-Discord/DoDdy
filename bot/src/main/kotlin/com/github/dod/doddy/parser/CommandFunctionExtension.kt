package com.github.dod.doddy.parser

import com.github.dod.doddy.core.CommandFunction
import net.dv8tion.jda.core.entities.Member
import net.dv8tion.jda.core.entities.User
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import java.lang.Exception
import kotlin.reflect.jvm.javaType


private val snowflakeRegex = Regex("\\d{17,19}")
private val mentionRegex = Regex("<@!?\\d{17,19}>")
private val usernameDiscrimRegex = Regex(".+#\\d{4}")

private val stringType = String::class.java
private val intType = Int::class.java
private val longType = Long::class.java
private val shortType = Short::class.java
private val doubleType = Double::class.java
private val userType = User::class.java
private val memberType = Member::class.java

suspend fun CommandFunction.call(event: MessageReceivedEvent, args: List<String>): CommandResult {
    if (args.size + optionals.size < parameters.size && !allArgs) {//TODO: check for too many arguments
        return InvalidArgs(args)
    }
    val params = ArrayList<Any?>()
    params.add(feature)
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
                userType -> {
                    val userId = when {
                        snowflakeRegex.matches(argument) -> argument
                        mentionRegex.matches(argument) -> argument.replace("<@!?".toRegex(), "").dropLast(1)
                        else -> null
                    }
                    val user: User? = when {
                        userId != null -> {
                            try {
                                event.jda.retrieveUserById(userId).complete()
                            } catch (_: Exception) {
                                null
                            }
                        }
                        usernameDiscrimRegex.matches(argument) -> {
                            val hashIndex = argument.lastIndexOf("#")
                            val username = argument.slice(0 until hashIndex)
                            val discrim = argument.substring(hashIndex + 1)
                            event.jda.getUsersByName(username, true).firstOrNull {
                                it.discriminator == discrim
                            }
                        }
                        else -> event.guild.getMembersByNickname(argument, true).firstOrNull()?.user
                            ?: event.jda.getUsersByName(argument, true).firstOrNull()
                    }
                    if (user != null) {
                        params.add(paramIndex, user)
                    } else {
                        return InvalidArg(argument, "not a valid user")
                    }
                }
                memberType -> {
                    val member: Member? = when {
                        snowflakeRegex.matches(argument) -> event.guild.getMemberById(argument)
                        mentionRegex.matches(argument) -> event.guild.getMemberById(argument.replace("^<@!?".toRegex(), "").dropLast(1))
                        usernameDiscrimRegex.matches(argument) -> {
                            val hashIndex = argument.lastIndexOf("#")
                            val username = argument.slice(0 until hashIndex)
                            val discrim = argument.substring(hashIndex + 1)
                            event.guild.getMembersByName(username, true).firstOrNull {
                                it.user.discriminator == discrim
                            }
                        }
                        else -> event.guild.getMembersByNickname(argument, true).firstOrNull()
                            ?: event.guild.getMembersByName(argument, true).firstOrNull()
                    }
                    if (member != null) {
                        params.add(paramIndex, member)
                    } else {
                        return InvalidArg(argument, "not a known member in this guild")
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