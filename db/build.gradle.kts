plugins {
    kotlin("jvm")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

dependencies {
    implementation(Libs.stdlib)
    api(Libs.coroutines)
    implementation(Libs.mongodb) { exclude(group = "org.jetbrains.kotlin", module = "kotlin-stdlib") }
}