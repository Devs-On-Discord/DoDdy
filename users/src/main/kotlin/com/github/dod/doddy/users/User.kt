package com.github.dod.doddy.users

data class User(val id: String,
                val guildReputations: Map<String, Long>) {
    val globalReputation: Long
    get() = guildReputations.values.max() ?: 0
}