package org.racingdb.app.repositories

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.jetbrains.exposed.sql.Transaction
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.selectAll
import org.racingdb.app.models.*

private fun ResultRow.toCategory() = Category(
    id = this[CategoriesTable.id],
    name = this[CategoriesTable.name],
    description = this[CategoriesTable.description],
    isActive = this[CategoriesTable.is_active]
)

class CategoryRepository {
    suspend fun getAll(tx: Transaction): List<Category> {
        var result = emptyList<Category>()
        withContext(Dispatchers.IO) {
            result = CategoriesTable.selectAll().map { it.toCategory() }
        }
        
        return result
    }
}