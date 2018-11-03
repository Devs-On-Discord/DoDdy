object Versions {
    val jda = "3.8.1_439"
    val kotlin = "1.3.0"
    val mongodb = "3.8.3"
    val coroutines = "1.0.0"
}

object Libs {
    val discord_bot_sdk = "net.dv8tion:JDA:${Versions.jda}"
    val stdlib = "org.jetbrains.kotlin:kotlin-stdlib-jdk8:${Versions.kotlin}"
    val reflection = "org.jetbrains.kotlin:kotlin-reflect:${Versions.kotlin}"
    val coroutines = "org.jetbrains.kotlinx:kotlinx-coroutines-core:${Versions.coroutines}"
    val mongodb = "org.litote.kmongo:kmongo-async:${Versions.mongodb}"
    val mongodb_coroutines = "org.litote.kmongo:kmongo-coroutine:${Versions.mongodb}"
}
