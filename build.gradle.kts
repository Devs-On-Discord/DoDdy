import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    base
    kotlin("jvm").version("1.3.0").apply(false)
}

allprojects {

    group = "com.github.dod.doddy"

    version = "1.0"

    tasks {
        withType<KotlinCompile> {
            kotlinOptions.jvmTarget = "1.8"
        }
    }

    repositories {
        jcenter()
    }
}