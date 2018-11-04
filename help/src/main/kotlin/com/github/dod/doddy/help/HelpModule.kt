package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.CommandFunction
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.AccountType
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.MessageBuilder
import net.dv8tion.jda.core.entities.MessageEmbed
import net.dv8tion.jda.core.events.message.MessageReceivedEvent

class HelpModule : Module {
    private val generalHelpPages = mutableListOf<MessageEmbed>()
    private val detailedHelpPages = mutableMapOf<String, MessageEmbed>()

    private val generalHelpPageTemplate = EmbedBuilder()
            .setColor(16777215)
            .setTitle("DoDdy", "https://github.com/Devs-On-Discord/DoDdy")

    @Command(names = ["help", "h"],
            shortDescription = "Displays help for commands.",
            longDescription = "Displays the help strings associated with every command. " +
                    "If a command is provided as an argument, the help string for that command is displayed instead.",
            docUrl = "https://github.com/Devs-On-Discord/DoDdy/blob/develop/help/README.md#help")
    fun help(event: MessageReceivedEvent, commandName: String?) {
        val privateChannelFuture = event.author.openPrivateChannel().submit()

        val confirmationMsg = MessageBuilder()
                .append("The requested information has been sent to your private messages, ")
                .append(event.author.asMention)
                .append(".").build()

        if (commandName == null) {
            event.channel.sendMessage(confirmationMsg).queue()

            val privateChannel = privateChannelFuture.get()
            generalHelpPages.forEach {
                privateChannel.sendMessage(it).queue()
            }
        } else {
            val detailedHelpPage = detailedHelpPages[commandName]
            if (detailedHelpPage == null) {
                val commandDoesNotExistMsg = MessageBuilder()
                        .append("I regret to inform you that there is no command called \"")
                        .append(commandName)
                        .append("\", ")
                        .append(event.author.asMention)
                        .append(".").build()

                event.channel.sendMessage(commandDoesNotExistMsg).queue()
                privateChannelFuture.cancel(true)
            } else {
                event.channel.sendMessage(confirmationMsg).queue()
                privateChannelFuture.get().sendMessage(detailedHelpPage).queue()
            }
        }
    }

    // Could be split up in more methods for more structure
    override fun onCommandsReady(commandFunctions: List<CommandFunction>) {
        var generalHelpPageBuilderIt = EmbedBuilder(generalHelpPageTemplate)

        val helpPageErrorNames = mutableListOf<String>()
        val generalHelpPageBuilders = mutableListOf<EmbedBuilder>()

        for (commandFunction in commandFunctions) {
            val commandMetaInfo = commandFunction.commandAnnotation

            val helpEntryProducer = HelpEntryProducer(commandFunction)
            val shortHelpEntry = helpEntryProducer.shortHelpEntry()
            val detailedHelpEntry = helpEntryProducer.detailedHelpEntry()

            // Check whether new general help page is needed (due to size limit)
            if (isValidWithNewEntry(generalHelpPageBuilderIt, shortHelpEntry)) {
                generalHelpPageBuilderIt.addField(shortHelpEntry)
            } else {
                generalHelpPageBuilders.add(generalHelpPageBuilderIt)
                generalHelpPageBuilderIt = EmbedBuilder(generalHelpPageTemplate)

                if (isValidWithNewEntry(generalHelpPageBuilderIt, shortHelpEntry))
                    generalHelpPageBuilderIt.addField(shortHelpEntry)
                else
                    helpPageErrorNames.add("[G] " + commandMetaInfo.names[0])
            }

            // Add detailed help pages if they can be sent
            if (detailedHelpEntry.isSendable(AccountType.BOT)) {
                addDetailedHelpPage(commandMetaInfo.names, detailedHelpEntry)
            } else {
                helpPageErrorNames.add("[D] " + commandMetaInfo.names[0])
            }
        }
        // Add remaining general help page
        generalHelpPageBuilders.add(generalHelpPageBuilderIt)

        // Add page numbers to footer
        val generalHelpPageAmount = generalHelpPageBuilders.size
        for (i in 0 until generalHelpPageAmount) {
            generalHelpPages.add(i, generalHelpPageBuilders[i]
                    .setFooter("Page " + (i+1) + "/" + generalHelpPageAmount, null)
                    .build())
        }

        if (helpPageErrorNames.size == 0) return

        // Turn error strings into embeds
        val errorMsgBuilder = EmbedBuilder()
                .setColor(0xFF0000)
                .setTitle("Errors when creating help pages")
                .setDescription("The following help pages couldn't be created: **")
        helpPageErrorNames.forEach {
            errorMsgBuilder.appendDescription("\n")
            errorMsgBuilder.appendDescription(it)
        }
        errorMsgBuilder.appendDescription("**")
        generalHelpPages.add(errorMsgBuilder.build())
    }

    private fun isValidWithNewEntry(helpPage: EmbedBuilder,
                                    shortHelpEntry: MessageEmbed.Field? = null) = EmbedBuilder(helpPage)
            .addField(shortHelpEntry)
            .setFooter("Page 999/999", null).build()    // Extreme case
            .isSendable(AccountType.BOT)

    private fun addDetailedHelpPage(commandNames: Array<out String>, detailedHelpEntry: MessageEmbed) {
        commandNames.forEach { commandName ->
            detailedHelpPages[commandName] = detailedHelpEntry
        }
    }
}