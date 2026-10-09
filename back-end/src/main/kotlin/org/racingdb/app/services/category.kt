package org.racingdb.app.services

import org.racingdb.app.core.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.racingdb.app.models.Category
import org.racingdb.app.database.MainDatabaseFactory

class CategoryService {
    
    private val core = CategoryCore()
    private val dbFactory = MainDatabaseFactory()
    
    suspend fun getAll(): List<Category> {
        var result = emptyList<Category>()
        dbFactory.executeInTransaction {tx -> 
            result = this.core.getAll(tx)
        }
        return result
    }
}