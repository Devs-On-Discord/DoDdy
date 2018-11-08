package com.github.dod.doddy.util

import net.dv8tion.jda.core.AccountType
import net.dv8tion.jda.core.EmbedBuilder
import net.dv8tion.jda.core.entities.MessageEmbed
import java.lang.IllegalStateException

fun <T : MessageEmbed.Field> Collection<T>.toMessageEmbeds(embedTemplate: EmbedBuilder, pageFooter: Boolean): MutableList<EmbedBuilder> {
    val messages = mutableListOf<EmbedBuilder>()

    var messageEmbedIt = EmbedBuilder(embedTemplate)

    for (field in this) {
        if (!isValidWithNewField(messageEmbedIt, field)) {
            messages.add(messageEmbedIt)
            messageEmbedIt = EmbedBuilder(embedTemplate)

            if (!isValidWithNewField(messageEmbedIt, field)) {
                IllegalStateException("\n- Field couldn't be parsed to message:\n${field.name}\n\n${field.value}\n").printStackTrace()
                continue
            }
        }
        messageEmbedIt.addField(field)
    }
    messages.add(messageEmbedIt)

    if (pageFooter) {
        val pageCount = messages.size

        for (page in 1 .. pageCount) {
            messages[page - 1].setFooter("Page $page/$pageCount", null)
        }
    }
    return messages
}

private fun isValidWithNewField(helpPage: EmbedBuilder,
                                shortHelpEntry: MessageEmbed.Field? = null) = EmbedBuilder(helpPage)
    .addField(shortHelpEntry)
    .setFooter("Page 999/999", null).build()    // Extreme case
    .isSendable(AccountType.BOT)