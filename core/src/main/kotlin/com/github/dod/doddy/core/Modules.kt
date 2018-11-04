package com.github.dod.doddy.core

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.launch
import net.dv8tion.jda.core.events.message.MessageReceivedEvent
import java.util.*
import kotlin.coroutines.CoroutineContext

class Modules internal constructor() : CoroutineScope {
    private val commands = Commands()

    private val modules = LinkedList<Module>()

    private lateinit var job: Job

    override val coroutineContext: CoroutineContext
        get() = Dispatchers.Default + job

    fun add(module: Module) {
        modules.add(module)
        commands.register(module)
    }

    internal suspend fun onCommand(name: String, event: MessageReceivedEvent, args: List<String>): CommandResult {
        return commands.call(name, event, args)
    }

    internal fun commandsReady() {
        job = Job()
        launch {
            repeat(modules.size) { i ->
                modules.elementAtOrNull(i)?.let { module ->
                    launch {
                        module.onCommandsReady(commands.functions)
                    }
                }
            }
        }
    }

    internal fun destroy() {
        modules.forEach { it.onDestroy() }
        job.cancel()
    }
}