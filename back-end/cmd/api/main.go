package main

import (
	"backend/internal/article"
	"backend/internal/articletext"
	"backend/internal/car"
	"backend/internal/category"
	"backend/internal/chassis"
	"backend/internal/circuit"
	"backend/internal/driver"
	"backend/internal/engine"
	"backend/internal/race"
	"backend/internal/team"
	"backend/modules/betterlogger"
	"backend/modules/database"
	"net/http"
	"os"
)

func main() {
	err := betterlogger.Init("app.log")
	if err != nil {
		panic(err)
	}
	defer betterlogger.Close()

	log := betterlogger.GetInstance()

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

	circuitHandler := circuit.NewCircuitHandler(db)
	circuitHandler.RegisterRoutes(mux)

	driverHandler := driver.NewDriverHandler(db)
	driverHandler.RegisterRoutes(mux)

	raceHandler := race.NewRaceHandler(db)
	raceHandler.RegisterRoutes(mux)

	teamHandler := team.NewTeamHandler(db)
	teamHandler.RegisterRoutes(mux)

	engineHandler := engine.NewEngineHandler(db)
	engineHandler.RegisterRoutes(mux)

	chassisHandler := chassis.NewChassisHandler(db)
	chassisHandler.RegisterRoutes(mux)

	carHandler := car.NewCarHandler(db)
	carHandler.RegisterRoutes(mux)

	articleHandler := article.NewArticleHandler(db)
	articleHandler.RegisterRoutes(mux)

	articletextHandler := articletext.NewArticleTextHandler(db)
	articletextHandler.RegisterRoutes(mux)

	port := ":5000"
	log.Info("Starting server on port %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Info("Server stopped: %v", err)
		os.Exit(0)
	}
}
