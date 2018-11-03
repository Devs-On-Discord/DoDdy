plugins {
    kotlin("jvm")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

dependencies {
    implementation(project(":db"))
    implementation(project(":guilds"))
    implementation(project(":users"))
    implementation("org.jetbrains.kotlin:kotlin-stdlib")
}