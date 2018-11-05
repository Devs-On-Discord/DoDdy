package com.github.dod.doddy.db

import com.mongodb.client.model.Updates
import org.bson.conversions.Bson
import org.litote.kmongo.path
import kotlin.reflect.KProperty

// Whoever this is, you're a blessing: https://stackoverflow.com/questions/25090548/push-equivalent-for-map-in-mongo
fun incEntry(property: KProperty<Map<*, *>>, entryName: String, number: Number): Bson = Updates.inc(property.path() + "." + entryName, number)
