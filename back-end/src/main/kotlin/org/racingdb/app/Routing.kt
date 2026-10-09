package org.racingdb.app

import io.ktor.server.application.*
import io.ktor.server.plugins.contentnegotiation.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import org.racingdb.app.handlers.*
    
fun Application.configureRouting() {
    install(ContentNegotiation) {
        json()
    }
    
    routing {
        miscRoutes()
        categoryRoutes()
    }
}