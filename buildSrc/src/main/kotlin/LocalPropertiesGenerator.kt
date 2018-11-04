import com.squareup.kotlinpoet.*
import groovy.lang.Closure
import org.gradle.api.DefaultTask
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.file.FileCollection
import org.gradle.api.internal.HasConvention
import org.gradle.api.plugins.JavaBasePlugin
import org.gradle.api.plugins.JavaPluginConvention
import org.gradle.api.tasks.*
import org.gradle.plugins.ide.idea.model.IdeaModel
import java.io.File
import java.util.concurrent.Callable

object LocalPropertiesGenerator {
    fun generate(file: File) {
        val greeterClass = ClassName("com.github.dod.doddy", "Greeter")
        val fileSpec = FileSpec.builder("com.github.dod.doddy", "HelloWorld")
                .addType(TypeSpec.classBuilder("Greeter")
                        .primaryConstructor(FunSpec.constructorBuilder()
                                .addParameter("name", String::class)
                                .build())
                        .addProperty(PropertySpec.builder("name", String::class)
                                .initializer("name")
                                .build())
                        .addFunction(FunSpec.builder("greet")
                                .addStatement("println(%S)", "Hello, \$name")
                                .build())
                        .build())
                .addFunction(FunSpec.builder("main")
                        .addParameter("args", String::class, KModifier.VARARG)
                        .addStatement("%T(args[0]).greet()", greeterClass)
                        .build())
                .build()
        fileSpec.writeTo(file)
    }
}

open class LocalPropertiesTask : DefaultTask() {

    private val DEFAULT_GEN_DIR = "gen"

    @OutputDirectory
    var outputDir: Any = project.file(DEFAULT_GEN_DIR)

    var classpath: FileCollection? = null

    init {
        group = "generate"
        outputs.upToDateWhen { task -> false } // do something smarter
    }

    @TaskAction
    fun generate() {
        LocalPropertiesGenerator.generate(project.file(outputDir))
    }
}

val Project.sourceSets: SourceSetContainer
    get() = project.convention.getPlugin(JavaPluginConvention::class.java).sourceSets

class CodegenPlugin : Plugin<Project> {

    override fun apply(project: Project) {
        project.plugins.apply(JavaBasePlugin::class.java)
        val sourceSetsToSkip = mutableSetOf("generators")
        project.sourceSets.forEach { sourceSet ->
            checkSourceSet(project, sourceSet, sourceSetsToSkip)
        }
    }

    private fun checkSourceSet(project: Project, sourceSet: SourceSet, sourceSetsToSkip: MutableSet<String>) {
        if (!sourceSetsToSkip.contains(sourceSet.name)) {
            if (sourceSet.name.endsWith("enerators")) {
                project.logger.error("Generators sourceSet (${sourceSet.name}) not registered in $sourceSetsToSkip")
            } else {
                processSourceSet(project, sourceSet, sourceSetsToSkip)
                project.logger.debug("sourceSetsToSkip is now: $sourceSetsToSkip")
            }
        }
    }

    private fun processSourceSet(project: Project, sourceSet: SourceSet, doSkip: MutableSet<String>) {
        val generateTaskName = if (sourceSet.name == "main") "generate" else sourceSet.getTaskName("generate", null)

        val generateConfiguration = project.configurations.maybeCreate(generateTaskName)
        project.configurations.add(generateConfiguration)

        val generatorSourceSetName = if (sourceSet.name == "main") "generators" else "${sourceSet.name}Generators"
        doSkip.add(generatorSourceSetName)

        val generatorSourceSet = project.sourceSets.maybeCreate(generatorSourceSetName)


        val outputDir = project.file("gen/${sourceSet.name}")

        val generateTask = project.tasks.create(generateTaskName, LocalPropertiesTask::class.java).apply {
            dependsOn(Callable { generateConfiguration })
            dependsOn(Callable { generatorSourceSet.classesTaskName })
            classpath = project.files(Callable { generateConfiguration }, Callable { generatorSourceSet.runtimeClasspath })
            this.outputDir = outputDir
        }

        project.dependencies.add(sourceSet.compileConfigurationName, project.files(Callable { generateConfiguration.files }).apply { builtBy(generateTask) })

        // Late bind the actual output directory
        sourceSet.java.srcDir(Callable { generateTask.outputDir })

        project.configurations.getByName(sourceSet.compileConfigurationName).extendsFrom(generateConfiguration)

        project.afterEvaluate {
            generateConfiguration.files(closure { generateTask.outputDir })

            project.extensions.findByType(IdeaModel::class.java)?.let { ideaModel ->
                ideaModel.module.generatedSourceDirs.add(project.file(generateTask.outputDir))
            }

            (project.tasks.getByName("clean") as? Delete)?.let { cleanTask ->
                cleanTask.delete(outputDir)
            }

        }
    }
}

private fun <T> closure(block: (args: Array<out Any?>) -> T): Closure<T> = object : Closure<T>(Unit) {

    override fun call(vararg args: Any?): T {
        return block(args)
    }
}