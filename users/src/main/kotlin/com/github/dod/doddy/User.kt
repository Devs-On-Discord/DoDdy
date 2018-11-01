package com.github.dod.doddy

data class User(val id: String,
                val guildReputations: Map<String, Long>) {
    val globalReputation: Long // Global reputation
    get() = guildReputations.values.sum()
}