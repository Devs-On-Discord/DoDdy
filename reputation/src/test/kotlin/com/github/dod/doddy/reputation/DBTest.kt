package com.github.dod.doddy.reputation

import com.github.dod.doddy.addReputation
import com.github.dod.doddy.users.User
import kotlinx.coroutines.runBlocking

object DBTest {
    @JvmStatic
    fun main(args: Array<String>) {
        val demo1 = User("snowflake01", HashMap())
        val demo2 = User("snowflake02", HashMap())
        runBlocking {
            demo1.addReputation("guildSnowflake01",
                    10)
            demo1.addReputation("guildSnowflake02",
                100)
            demo1.addReputation("guildSnowflake01",
                10)
            demo1.addReputation("guildSnowflake01",
                10)
            demo1.addReputation("guildSnowflake02",
                50)
            demo1.addReputation("guildSnowflake01",
                10)
            demo2.addReputation("guildSnowflake02",
                1337)
        }
    }
}