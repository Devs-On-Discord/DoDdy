package com.github.dod.doddy.users

data class User(val id: String,
                val guildReputations: Map<String, Long>) {
    fun calcGlobalReputation() = guildReputations.values.max() ?: 0
}