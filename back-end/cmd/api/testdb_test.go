package main

import (
	"backend/modules/database"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// The integration tests run the real handler against a Postgres loaded with ../../../init.sql and
// ../../../seed.sql. The server comes from TEST_DATABASE_URL when set (e.g. the compose one), else from a throwaway
// testcontainers container. Either way the tests get a fresh database of their own, dropped afterwards.
// With no URL and no Docker the tests are skipped; with a URL set, any setup error fails the run.

var (
	testDB      *sql.DB
	testHandler http.Handler
	skipReason  string
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	serverURL, stopServer, err := testServer(ctx)
	if err != nil {
		if os.Getenv("TEST_DATABASE_URL") != "" {
			fmt.Fprintln(os.Stderr, "integration tests:", err)
			return 1
		}
		skipReason = fmt.Sprintf("no TEST_DATABASE_URL and no Docker (%v)", err)
		return m.Run()
	}
	defer stopServer()

	db, dropDB, err := freshDatabase(ctx, serverURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration tests:", err)
		return 1
	}
	defer dropDB()

	// Handlers log client errors through the default logger; the test output only needs failures.
	slog.SetDefault(slog.New(slog.DiscardHandler))
	testDB = db
	testHandler = newHandler(slog.Default(), db)
	return m.Run()
}

// requireDB skips the calling test when there is no database to test against.
func requireDB(t *testing.T) {
	t.Helper()
	if testDB == nil {
		t.Skip("integration test skipped: " + skipReason)
	}
}

// testServer returns a connection string for a Postgres server the tests may create databases on.
func testServer(ctx context.Context) (string, func(), error) {
	if u := os.Getenv("TEST_DATABASE_URL"); u != "" {
		return u, func() {}, nil
	}

	ctr, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("racingdb"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("adminpassword"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return "", nil, err
	}
	stop := func() { _ = ctr.Terminate(context.Background()) }

	u, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		stop()
		return "", nil, err
	}
	return u, stop, nil
}

// freshDatabase creates a uniquely named database on the server, loads the schema and seed into it, and
// returns a pool on it plus a function that closes the pool and drops the database.
func freshDatabase(ctx context.Context, serverURL string) (*sql.DB, func(), error) {
	admin, err := database.NewPostgresDB(serverURL)
	if err != nil {
		return nil, nil, err
	}

	name := fmt.Sprintf("racingdb_test_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+name); err != nil {
		admin.Close()
		return nil, nil, fmt.Errorf("create database: %w", err)
	}

	var db *sql.DB
	drop := func() {
		if db != nil {
			db.Close()
		}
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
		admin.Close()
	}

	u, err := url.Parse(serverURL)
	if err != nil {
		drop()
		return nil, nil, err
	}
	u.Path = "/" + name

	if db, err = database.NewPostgresDB(u.String()); err != nil {
		drop()
		return nil, nil, err
	}

	// Executed without arguments, so pgx sends each file as one multi-statement simple query.
	for _, file := range []string{"init.sql", "seed.sql"} {
		script, err := os.ReadFile(filepath.Join("..", "..", "..", file))
		if err != nil {
			drop()
			return nil, nil, err
		}
		if _, err := db.ExecContext(ctx, string(script)); err != nil {
			drop()
			return nil, nil, fmt.Errorf("load %s: %w", file, err)
		}
	}

	return db, drop, nil
}
