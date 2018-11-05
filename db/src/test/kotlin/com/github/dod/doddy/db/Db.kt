package com.github.dod.doddy.db

import com.github.dod.doddy.db.coroutines.insertOne
import kotlinx.coroutines.runBlocking
import org.litote.kmongo.async.KMongo
import org.litote.kmongo.async.getCollection

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