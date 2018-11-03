plugins {
    kotlin("jvm")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

dependencies {
    implementation(Libs.stdlib)
    implementation(Libs.discord_bot_sdk)
    implementation(Libs.reflection)
}