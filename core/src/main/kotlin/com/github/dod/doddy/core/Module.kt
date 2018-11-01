package com.github.dod.doddy.core

interface Module {
    /**
     * @return true to consume
     */
    fun onCommand(command: String, args: Array<Any>): Boolean

    fun getCommands(): Set<String>

    fun getEventListeners() : List<Any>
}