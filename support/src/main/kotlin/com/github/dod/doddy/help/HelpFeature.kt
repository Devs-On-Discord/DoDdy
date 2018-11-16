package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.CommandFunction
import com.github.dod.doddy.core.Feature
import com.github.dod.doddy.util.toMessageEmbeds
import net.dv8tion.jda.core.AccountType
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.JDA
import net.dv8tion.jda.core.MessageBuilder
import net.dv8tion.jda.core.entities.MessageEmbed
import net.dv8tion.jda.core.events.message.MessageReceivedEvent

class HelpFeature : Feature {
    private val generalHelpPages = mutableListOf<MessageEmbed>()
    private val detailedHelpPages = mutableMapOf<String, MessageEmbed>()

    private val generalHelpPageTemplate = EmbedBuilder()
            .setColor(0xFFFFFF)
            .setTitle("DoDdy", "https://github.com/Devs-On-Discord/DoDdy")

    @Command(names = ["help", "h"],
            shortDescription = "Displays help for commands.",
            longDescription = "Displays the help strings associated with every command. " +
                    "If a command is provided as an argument, the help string for that command is displayed instead.",
            docUrl = "https://github.com/Devs-On-Discord/DoDdy/tree/develop/support#help")
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
    override fun onBotReady(bot: JDA, commandFunctions: List<CommandFunction>) {
        val generalHelpPageFields = mutableListOf<MessageEmbed.Field>()

        for (commandFunction in commandFunctions) {
            val commandMetaInfo = commandFunction.commandAnnotation

            val helpEntryProducer = HelpEntryProducer(commandFunction)
            val shortHelpEntry = helpEntryProducer.shortHelpEntry()
            val detailedHelpEntry = helpEntryProducer.detailedHelpEntry()

            generalHelpPageFields.add(shortHelpEntry)

            // Add detailed help pages if they can be sent
            if (detailedHelpEntry.isSendable(AccountType.BOT)) {
                commandMetaInfo.names.forEach { commandName ->
                    detailedHelpPages[commandName] = detailedHelpEntry
                }
            } else {
                IllegalStateException("Couldn't render detailed help page for command: ${commandMetaInfo.names[0]}").printStackTrace()
            }
        }

        generalHelpPageFields.toMessageEmbeds(generalHelpPageTemplate, true).forEach {
            generalHelpPages.add(it.build())
        }
    }
}