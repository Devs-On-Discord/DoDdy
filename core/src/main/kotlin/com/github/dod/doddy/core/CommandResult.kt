package com.github.dod.doddy.core

sealed class CommandResult
data class CommandNotFound(val commandName: String) : CommandResult()
data class InvalidArgs(val args: List<String>) : CommandResult()
data class InvalidArg(val arg: String, val error: String) : CommandResult()
data class CommandError(val exception: Exception) : CommandResult()
data class Success(val message: String) : CommandResult()