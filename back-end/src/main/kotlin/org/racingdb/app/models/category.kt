package org.racingdb.app.models

import org.jetbrains.exposed.sql.Table
import kotlinx.serialization.Serializable

object CategoriesTable : Table("categories") {
    val id = varchar("id", 50)
    val name = varchar("name", 100)
    val description = text("description")
    val is_active = bool("is_active")
    
    override val primaryKey = PrimaryKey(id)
}

@Serializable
data class Category(
    val id: String,
    val name: String,
    val description: String,
    val isActive: Boolean
)