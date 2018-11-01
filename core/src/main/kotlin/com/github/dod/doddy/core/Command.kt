package com.github.dod.doddy.core

/**
 * Parent class for all commands used by the Discord bot.
 *
 * @property names the names of the command. Sorted by frequency of intended use.
 * Example: ["reputation", "rep", "re"]
 * "reputation" would be used as the command main identifier, f.e. the "help" command.
 */
abstract class Command(val names: Array<String>) {
    /**
     * Called when the bot detects usage of one of the specified names (constructor argument)
     *
     * @return true to consume
     */
    abstract fun onCommand(args: Array<Any>): Boolean
}