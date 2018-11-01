package com.github.dod.doddy.help

import com.github.dod.doddy.core.Module

class HelpModule : Module {
    override fun getCommands() = setOf(
            HelpCommand())

    override fun getEventListeners(): List<Any> {
        return emptyList()
    }
}