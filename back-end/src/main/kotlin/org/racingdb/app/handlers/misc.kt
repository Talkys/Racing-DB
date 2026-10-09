package org.racingdb.app.handlers

import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
    
fun Route.miscRoutes() {
    route("/") {
        get {
            call.respondText("Hello, World!")
        }
    }
}