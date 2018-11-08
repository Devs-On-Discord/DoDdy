plugins {
    application
    kotlin("jvm")
    CodegenPlugin()
}

application {
    mainClassName = "com.github.dod.doddy.Main"
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

//TODO: this doesn't work
/*tasks {
    val localProperties by registering(LocalPropertiesTask::class)

    localProperties {
        dependsOn(build)
    }
}*/

dependencies {
    implementation(project(":core"))

    implementation(project(":support"))
    implementation(project(":moderation"))
    implementation(project(":info"))
    implementation(project(":reputation"))
    implementation(project(":polls"))
    implementation(Libs.stdlib)
    implementation(Libs.coroutines)
    implementation(Libs.discord_bot_sdk)
    implementation(Libs.reflection)
}