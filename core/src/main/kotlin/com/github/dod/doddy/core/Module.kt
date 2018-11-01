package com.github.dod.doddy.core

interface Module {
    /**
     * @return true to forward commands, false to consume
     */
    fun onCommand(): Boolean
}