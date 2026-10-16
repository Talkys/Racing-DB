package main

import (
	"backend/internal/category"
	logger "backend/modules/betterlogger"
	"backend/modules/database"
	"net/http"
	"os"
)

func main() {
	err := logger.Init("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	log := logger.GetInstance()

	log.Info("Starting application initialization...")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Warn("DATABASE_URL environment variable not found. Using local default connection string.")
		connStr = "postgres://admin:adminpassword@localhost:5432/racingdb?sslmode=disable"
	}

	log.Info("Connecting to PostgreSQL database pool...")
	db, err := database.NewPostgresDB(connStr)
	if err != nil {
		log.Error("Initialization error: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Success("Database connection pool established successfully.")

	mux := http.NewServeMux()

	log.Info("Registering routes...")
	// REGISTER HANDLERS HERE (MAYBE THERE'S A WAY TO DO BETTER)
	categoryHandler := category.NewCategoryHandler(db)
	categoryHandler.RegisterRoutes(mux)

	port := ":5000"
	log.Info("Starting server on port %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Info("Server stopped: %v", err)
		os.Exit(0)
	}
}
