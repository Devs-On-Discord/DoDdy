package com.github.dod.doddy.reputation

import com.github.dod.doddy.addReputation
import com.github.dod.doddy.db.Db
import com.github.dod.doddy.db.coroutines.dropCollection
import com.github.dod.doddy.users.User
import kotlinx.coroutines.runBlocking

object DBTest {
    @JvmStatic
    fun main(args: Array<String>) {
        // Mongo cache broken? Users always created correctly
        for (x in 0 until 50) {
            val demo1 = User("snowflake01", HashMap())
            val demo2 = User("snowflake02", HashMap())
            println(demo1)
            println(demo2)
            runBlocking {
                Db.instance.dropCollection("user")
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
            println("-!-")
        }
    }
}