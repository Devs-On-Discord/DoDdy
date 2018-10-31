plugins {
    application
    kotlin("jvm")
}

application {
    mainClassName = "com.github.doddy.Main"
}

dependencies {
    implementation(project(":feature1"))
    implementation("org.jetbrains.kotlin:kotlin-stdlib")
}