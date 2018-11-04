package com.github.dod.doddy.reputation

import com.github.dod.doddy.addReputation
import com.github.dod.doddy.users.User
import kotlinx.coroutines.runBlocking

object DBTest {
    @JvmStatic
    fun main(args: Array<String>) {
        runBlocking {
            val demo1 = User("snowflake01", emptyMap())
            demo1.addReputation("guildSnowflake01",
                    11)
            demo1.addReputation("guildSnowflake02",
                15)
            demo1.addReputation("guildSnowflake01",
                14)
            val demo2 = User("snowflake02", emptyMap())
            demo2.addReputation("guildSnowflake02",
                1337)
        }
    }
}