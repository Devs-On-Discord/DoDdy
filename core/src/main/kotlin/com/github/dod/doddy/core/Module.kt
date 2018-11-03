package com.github.dod.doddy.core

interface Module {
    /**
     * @return the event listeners that should be added to the discord bot
     */
    fun getEventListeners() : List<Any>
}