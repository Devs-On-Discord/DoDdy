@file:JvmName("Main")
package com.github.dod.doddy

import com.github.dod.doddy.parser.Discord
import com.github.dod.doddy.bans.BansFeature
import com.github.dod.doddy.help.HelpFeature

fun main(args : Array<String>) {
    println("Hello, world!")
    Discord.add(BansFeature())
    Discord.add(HelpFeature())
    Discord.start()
}