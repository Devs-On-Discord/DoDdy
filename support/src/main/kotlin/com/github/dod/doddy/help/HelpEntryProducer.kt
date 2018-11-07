package com.github.dod.doddy.help

import com.github.dod.doddy.core.CommandFunction
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.entities.MessageEmbed

class HelpEntryProducer(private val commandFunction: CommandFunction) {

    fun shortHelpEntry(): MessageEmbed.Field {
        val commandMetaInfo = commandFunction.commandAnnotation
        return MessageEmbed.Field(commandMetaInfo.names[0], commandMetaInfo.shortDescription,false)
    }

    fun detailedHelpEntry(): MessageEmbed {
        val embed = EmbedBuilder()
                .setTitle(generateCommandUsage(), commandFunction.commandAnnotation.docUrl)
                .setDescription("**")
                .appendDescription(commandFunction.commandAnnotation.shortDescription)
                .appendDescription("**\n\n")
                .appendDescription(commandFunction.commandAnnotation.longDescription)
                .setColor(0xFFFFFF)

        val commandNames = commandFunction.commandAnnotation.names
        if (commandNames.size > 1) {
            val footerBuilder = StringBuilder("Aliases:")
            commandFunction.commandAnnotation.names.forEach { footerBuilder.append(" $it,") }
            embed.setFooter(footerBuilder.dropLast(1).toString(), null)
        }

        return embed.build()
    }

    private fun generateCommandUsage(): String {

        val command = commandFunction.commandAnnotation
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