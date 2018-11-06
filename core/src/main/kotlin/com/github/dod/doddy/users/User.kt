package com.github.dod.doddy.users

data class User(val id: String,
                val guildReputations: Map<String, Long>) {
    fun getGlobalReputation() = guildReputations.values.max() ?: 0
}