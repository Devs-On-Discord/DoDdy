package com.github.dod.doddy.core

interface Module {
    /**
     * @return the commands that should exclusively executed by this module
     */
    fun getCommands(): Set<Command>

    /**
     * @return the event listeners that should be added to the discord bot
     */
    fun getEventListeners() : List<Any>
}