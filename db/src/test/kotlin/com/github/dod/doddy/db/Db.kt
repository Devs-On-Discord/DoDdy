package com.github.dod.doddy.db

import kotlinx.coroutines.runBlocking
import org.litote.kmongo.async.KMongo
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.coroutine.insertOne

object Db {
    val client = KMongo.createClient()
    val database = client.getDatabase("mongo")

    @JvmStatic
    fun main(args: Array<String>) {
        main()
    }

    @JvmStatic
    fun main() = runBlocking {
        val col = database.getCollection<Car>()
        col.insertOne(Car(2, "Brumma"))
    }
}