package com.github.dod.doddy.core

/**
 * Parent class for all commands used by the Discord bot.
 *
 * @property names the names of the command. Sorted by frequency of intended use.
 * Example: ["reputation", "rep", "re"]
 * "reputation" would be used as the command main identifier, f.e. used by the "help" command.
 */
abstract class Command(vararg val names: String) {
    /**
     * Called when the bot detects usage of one of the specified *names*
     *
     * @return true to consume
     */
    abstract fun onCommand(args: List<String>): Boolean
}