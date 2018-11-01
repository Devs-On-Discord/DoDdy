plugins {
    kotlin("jvm")
}

dependencies {
    implementation(project(":db"))
    implementation(project(":guilds"))
    implementation(project(":users"))
    implementation(Libs.stdlib)
}