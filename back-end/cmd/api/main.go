package main

import (
	"backend/internal/article"
	"backend/internal/articlesearch"
	"backend/internal/articletext"
	"backend/internal/car"
	"backend/internal/category"
	"backend/internal/chassis"
	"backend/internal/circuit"
	"backend/internal/crud"
	"backend/internal/driver"
	"backend/internal/engine"
	"backend/internal/race"
	"backend/internal/team"
	"backend/modules/database"
	"backend/modules/logging"
	"backend/utils/constants"
	"backend/utils/routes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title RacingDB API
// @version 1
// @description Read-only racing database: categories, circuits, drivers, races, teams, engines, chassis, cars and articles.
// @BasePath /api/1
func main() {
	log, err := logging.New(os.Stdout, os.Getenv("LOG_FORMAT"), os.Getenv("LOG_LEVEL"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	slog.SetDefault(log)

	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

// run holds the whole lifetime of the server, so its defers (db.Close) run before main exits.
func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Warn("DATABASE_URL not set, using the local default connection string")
		connStr = "postgres://admin:adminpassword@localhost:5432/racingdb?sslmode=disable"
	}

	db, err := database.NewPostgresDB(connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("database connection pool established")

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":5000"
	}

	srv := &http.Server{
		Addr:              port,
		Handler:           newHandler(log, db),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Info("server listening", "addr", port)
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}
	stop() // a second Ctrl+C now kills the process instead of waiting

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Info("server stopped")
	return nil
}

// newHandler wires every API route on one mux, wrapped in the request middleware.
func newHandler(log *slog.Logger, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	handlers := []interface{ RegisterRoutes(*http.ServeMux) }{
		crud.NewHandler(db, category.Resource),
		crud.NewHandler(db, circuit.Resource),
		crud.NewHandler(db, driver.Resource),
		crud.NewHandler(db, race.Resource),
		crud.NewHandler(db, team.Resource),
		crud.NewHandler(db, engine.Resource),
		crud.NewHandler(db, chassis.Resource),
		crud.NewHandler(db, car.Resource),
		crud.NewHandler(db, article.Resource),
		crud.NewHandler(db, articletext.Resource),
		articlesearch.NewArticleSearchHandler(db),
	}
	for _, h := range handlers {
		h.RegisterRoutes(mux)
	}

	// Anything under the API prefix that no handler claimed gets a JSON 404 instead of the mux's plain text.
	mux.HandleFunc(constants.BaseApiUrl_V1+"/", routes.NotFound)

	return routes.Middleware(log, mux)
}
