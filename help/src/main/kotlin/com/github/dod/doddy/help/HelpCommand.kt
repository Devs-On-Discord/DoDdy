package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command

class HelpCommand : Command(arrayOf("help", "h")) {
    override fun onCommand(args: Array<Any>): Boolean {
        return false
    }
}