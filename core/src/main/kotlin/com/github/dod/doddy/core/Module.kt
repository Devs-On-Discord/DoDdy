package com.github.dod.doddy.core

interface Module {
    /**
     * @return true to consume
     */
    fun onCommand(command: String, args: Array<Any>): Boolean

    /**
     * @return the commands that should exclusively executed by this module
     */
    fun getCommands(): Set<String>

    /**
     * @return the event listeners that should be added to the discord bot
     */
    fun getEventListeners() : List<Any>
}