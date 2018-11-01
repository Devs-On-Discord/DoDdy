plugins {
    kotlin("jvm")
}

dependencies {
    implementation(project(":db"))
    implementation(project(":guilds"))
    implementation(project(":users"))
    implementation("org.jetbrains.kotlin:kotlin-stdlib")
}