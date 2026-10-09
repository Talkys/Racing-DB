package org.racingdb.app.database

import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.Transaction
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction

class MainDatabaseFactory {
    val db: Database

    init {
        // TODO: Colocar dados de acesso em env
        val config = HikariConfig().apply {
            this.jdbcUrl = "jdbc:postgresql://localhost:5432/racingdb"
            driverClassName = "org.postgresql.Driver"
            username = "admin"
            password = "adminpassword"
            maximumPoolSize = 10
            isAutoCommit = false
            transactionIsolation = "TRANSACTION_REPEATABLE_READ"
            validate()
        }
        db = Database.connect(HikariDataSource(config))
    }

    suspend fun <T> executeInTransaction(block: suspend (Transaction) -> T): T {
        return newSuspendedTransaction {
            block(this)
        }
    }
}