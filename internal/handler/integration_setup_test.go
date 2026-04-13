package handler_test

import (
	"context"
	"database/sql"
	"mini-ledger/internal/handler"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName         = "dev"
	dbUser         = "dev"
	dbPasswd       = "dev"
	postgresImage  = "postgres:17-alpine"
	migrationsPath = "file://../../migrations"
)

func runAllMigrations(t *testing.T, ctx context.Context, db *sql.DB) *migrate.Migrate {
	t.Helper()
	conn, err := db.Conn(ctx)
	require.NoError(t, err)

	driver, err := migratePostgres.WithConnection(
		ctx,
		conn,
		&migratePostgres.Config{}, //nolint:exhaustruct
	)
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"pgx",
		driver)

	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	return m
}

func allocateResources(t *testing.T) (humatest.TestAPI, func()) {
	t.Helper()
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		postgresImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPasswd),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			).WithStartupTimeoutDefault(30*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	_, api := humatest.New(t)
	handlers := handler.WireLayers(db)

	handler.RegisterAllRoutes(handlers, api)

	m := runAllMigrations(t, ctx, db)

	return api, func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Error().Err(sourceErr).Msg("Failed to close migration source")
		}
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Failed to close migration database connection")
		}
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate postgres container: %v", err)
		}
	}
}
