package com.github.dod.doddy.help

import com.github.dod.doddy.core.CommandFunction
import net.dv8tion.jda.core.EmbedBuilder

class CommandDescriptionBuilder(private val commandFunction: CommandFunction) {
    fun shortDesc(pEmbedBuilder: EmbedBuilder): EmbedBuilder {
        val embedBuilder = EmbedBuilder(pEmbedBuilder)

        embedBuilder.addField(generateCommandUsage(), commandFunction.command.shortDescription, false)
        return embedBuilder
    }

    fun longDesc(embedBuilder: EmbedBuilder) {
        embedBuilder.setTitle(generateCommandUsage(), commandFunction.command.docUrl)

        embedBuilder.setDescription("**")
        embedBuilder.appendDescription(commandFunction.command.shortDescription)
        embedBuilder.appendDescription("**\n\n")
        embedBuilder.appendDescription(commandFunction.command.longDescription)

        embedBuilder.setColor(16777215)

        val footerBuilder = StringBuilder("Aliases:")
        commandFunction.command.names.forEach { footerBuilder.append(" $it,") }

        footerBuilder.dropLast(1)

        embedBuilder.setFooter(footerBuilder.toString(), null)
    }

    private fun generateCommandUsage(): String {
        val command = commandFunction.command
        val commandUsageBuilder = StringBuilder(command.names[0])

        for (i in 0 until commandFunction.parameters.size) {
            val parameter = commandFunction.parameters[i]

            commandUsageBuilder.append(" ")

            commandUsageBuilder.append(if (parameter.type.isMarkedNullable) "[" else "<")
            commandUsageBuilder.append(parameter.name)
            if (i == commandFunction.parameters.size - 1 && commandFunction.allArgs)
                commandUsageBuilder.append("...")
            commandUsageBuilder.append(if (parameter.type.isMarkedNullable) "]" else ">")
        }

        return commandUsageBuilder.toString()
    }
}