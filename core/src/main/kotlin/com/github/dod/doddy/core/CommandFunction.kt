package com.github.dod.doddy.core

import kotlin.reflect.KFunction
import kotlin.reflect.KParameter

data class CommandFunction(
    val feature: Feature,
    val function: KFunction<*>,
    val parameters: List<KParameter>,
    val optionals: List<Int>,
    val allArgs: Boolean,
    val commandAnnotation: Command
)