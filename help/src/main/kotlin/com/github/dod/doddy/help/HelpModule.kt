package com.github.dod.doddy.help

import com.github.dod.doddy.core.Command
import com.github.dod.doddy.core.CommandFunction
import com.github.dod.doddy.core.Module
import net.dv8tion.jda.core.AccountType
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.MessageBuilder
import net.dv8tion.jda.core.events.message.MessageReceivedEvent

// PROTOTYPE VERSION: Will be redone code-wise by ThisIsIPat, no worries about how it looks for now.
class HelpModule : Module {
    lateinit var commandFunctions: List<CommandFunction>

    @Command(names = ["help", "h"],
            shortDescription = "Displays help for commands.",
            longDescription = "Displays the help strings associated with every command." +
                    "If a command is provided as an argument, the help string for that command is displayed instead.",
            docUrl = "https://github.com/Devs-On-Discord/DoDdy/blob/develop/help/README.md#help")
    fun help(event: MessageReceivedEvent, commandName: String?) {

        val privateChannelFuture = event.author.openPrivateChannel().submit()
        val confirmationMessage = MessageBuilder()
                .append("The requested information has been sent to your private messages, ")
                .append(event.author.asMention)
                .append(".").build()

        var embedBuilder = EmbedBuilder()
        if (commandName != null) {
            val command = commandFunctions.firstOrNull { it.command.names.contains(commandName) }
            if (command != null) {
                CommandDescriptionBuilder(command).longDesc(embedBuilder)
                val embed = embedBuilder.build()

                if (embed.isSendable(AccountType.BOT)) {
                    event.channel.sendMessage(confirmationMessage).queue()
                    privateChannelFuture.get().sendMessage(embed).queue()
                } else {
                    event.channel.sendMessage(
                            "Unfortunately, the command $commandName provided an invalid help text.").queue()
                    privateChannelFuture.cancel(true)
                }
            } else {
                event.channel.sendMessage(MessageBuilder()
                        .append("I regret to inform you that there is no command called \"")
                        .append(commandName)
                        .append("\", ")
                        .append(event.author.asMention)
                        .append(".").build()).queue()
                privateChannelFuture.cancel(true)
            }
        } else {
            // Prepare multiple help pages
            val helpEmbedBuilders = mutableListOf<EmbedBuilder>()

            commandFunctions.forEach { commandFunction ->
                val commandDescriptionBuilder = CommandDescriptionBuilder(commandFunction)
                val unverifiedEmbedBuilder = commandDescriptionBuilder.shortDesc(embedBuilder)

                if (unverifiedEmbedBuilder.build().isSendable(AccountType.BOT)) {
                    embedBuilder = unverifiedEmbedBuilder
                } else {
                    helpEmbedBuilders.add(embedBuilder)
                    embedBuilder = commandDescriptionBuilder.shortDesc(EmbedBuilder())
                    if (!embedBuilder.build().isSendable(AccountType.BOT)) {
                        embedBuilder = EmbedBuilder().appendDescription(
                                "Unfortunately, the command " + commandFunction.command.names[0] +
                                        " provided an invalid help text.\n\n")
                    }
                }
            }
            helpEmbedBuilders.add(embedBuilder)

            event.channel.sendMessage(confirmationMessage).queue()
            val privateChannel = privateChannelFuture.get()
            for (i in 0 until helpEmbedBuilders.size) {
                helpEmbedBuilders[i].setTitle("DoDdy", "https://github.com/Devs-On-Discord/DoDdy")
                helpEmbedBuilders[i].setColor(16777215)
                helpEmbedBuilders[i].setFooter("Page " + (i+1) + "/" + helpEmbedBuilders.size, null)
                if (helpEmbedBuilders[i].build().isSendable(AccountType.BOT))  {
                    privateChannel.sendMessage(helpEmbedBuilders[i].build()).queue()
                } else {
                    privateChannel.sendMessage("Oops! There has been a problem sending a help page!").queue()
                }
            }
        }
    }

    // For a bored programmer:
    // Improve performance by one-time-calculating the help embeds when this method is called:
    override fun onCommandFunctionsReady(commandFunctions: List<CommandFunction>) {
        this.commandFunctions = commandFunctions
    }
}