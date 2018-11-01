package com.github.dod.doddy.help

import com.github.dod.doddy.core.Module
import java.util.*

class HelpModule: Module {
    override fun onCommand(command: String, args: Array<Any>): Boolean {

        return true
    }

    override fun getCommands(): Set<String> {
        return setOf("help")
    }

    override fun getEventListeners(): List<Any> {
        return Collections.emptyList()
    }
}