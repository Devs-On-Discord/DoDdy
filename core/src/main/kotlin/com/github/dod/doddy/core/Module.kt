package com.github.dod.doddy.core

interface Module {
    /**
     * @return true to consume
     */
    fun onCommand(): Boolean
}