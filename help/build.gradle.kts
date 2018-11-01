plugins {
    kotlin("jvm")
}

dependencies {
    implementation(project(":core"))
    implementation(Libs.stdlib)
    implementation(Libs.discord_bot_sdk)
}