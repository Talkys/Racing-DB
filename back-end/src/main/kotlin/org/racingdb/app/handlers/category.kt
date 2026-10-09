package org.racingdb.app.handlers

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import org.racingdb.app.services.*
    
fun Route.categoryRoutes() {
    route("/category") {
        get("/all/") {
            val service = CategoryService()
            val result = service.getAll()
            
            println(result)
            
            call.respond(HttpStatusCode.OK, result)
        }
    }
}