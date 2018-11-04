@file:JvmName("Main")
package com.github.dod.doddy

import com.github.dod.doddy.core.Discord
import com.github.dod.doddy.bans.BansModule
import com.github.dod.doddy.help.HelpModule

fun main(args : Array<String>) {
    println("Hello, world!")
    Discord.add(BansModule())
    Discord.add(HelpModule())
    Discord.start()
}