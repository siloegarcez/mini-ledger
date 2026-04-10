package main

import (
	"context"
	"database/sql"
	"errors"
	"mini-ledger/config"
	"mini-ledger/database"
	"mini-ledger/internal/server"
	"mini-ledger/logging"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, conf *config.Config) {
		logging.Configure(*conf)

		errs := conf.Validate()
		if len(errs) != 0 {
			log.Fatal().Errs("errs", errs).Msg("Invalid configuration, stopping server startup")
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(conf.StartupTimeoutSec)*time.Second,
		)
		defer cancel()

		db, err := database.Connect(ctx, *conf)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to connect to the database, stopping server startup")
		}

		if conf.Environment == config.EnvDev {
			log.Info().Msg("Running database migrations")
			err := migrateUp(ctx, db, *conf)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("Failed to run database migrations, Stopping server startup")
			} else {
				log.Info().Msg("Database migrations completed successfully")
			}
		}

		srv := server.New(*conf, db)

		hooks.OnStart(func() {
			srv.RegisterRoutes()

			log.Info().
				Msgf("%s - (v%s) running on port %d", conf.Name, conf.Version, conf.Port)
			err = srv.ListenAndServe()

			if errors.Is(err, http.ErrServerClosed) {
				log.Info().Msg("Server closed")
			} else {
				log.Error().Err(err).Msg("Server startup failed")
			}
		})

		hooks.OnStop(func() {
			log.Info().
				Msgf("Shutting down the server, all goroutines must finish the current request within %d seconds", conf.GracefulShutdownTimeoutSec)

			err := srv.GracefulShutdown()

			if err != nil {
				log.Error().Err(err).Msg("Server shutdown encountered an error")
			} else {
				log.Info().Msg("Graceful shutdown complete")
			}
		})
	})

	cli.Run()
}

func migrateUp(ctx context.Context, db *sql.DB, conf config.Config) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}

	driver, err := postgres.WithConnection(ctx, conn, &postgres.Config{}) //nolint:exhaustruct

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+conf.MigrationsPath,
		"pgx",
		driver)

	if err != nil {
		return err
	}

	err = m.Up()
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Error().Err(sourceErr).Msg("Failed to close migration source")
		}
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Failed to close migration database connection")
		}
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
