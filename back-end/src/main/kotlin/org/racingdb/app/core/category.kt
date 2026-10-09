package org.racingdb.app.core

import org.racingdb.app.repositories.*
import org.racingdb.app.models.Category
import org.jetbrains.exposed.sql.Transaction
    
class CategoryCore {
    private val repository = CategoryRepository() 
    
    suspend fun getAll(tx: Transaction): List<Category> {
        return this.repository.getAll(tx)
    }
}