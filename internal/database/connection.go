package database

import (
	"context"
	"database/sql"
	"mini-ledger/config"

	"net/url"
	"time"
)

func Connect(ctx context.Context, conf config.Config) (*sql.DB, error) {
	u := &url.URL{ //nolint: exhaustruct
		Scheme: "postgres",
		Host:   conf.DBHost + ":" + conf.DBPort,
		Path:   conf.DBName,
	}

	u.User = url.UserPassword(conf.DBUsername, conf.DBPassword)

	q := u.Query()

	if conf.Environment == config.EnvProd || conf.Environment == config.EnvStaging {
		q.Set("sslmode", "require")
	} else {
		q.Set("sslmode", "disable")
	}

	if conf.DBSchema != "" {
		q.Set("search_path", conf.DBSchema)
	}

	u.RawQuery = q.Encode()

	connStr := u.String()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(conf.DBConnMaxLifetimeSec) * time.Second)
	db.SetMaxIdleConns(conf.DBMaxIdleConns)
	db.SetMaxOpenConns(conf.DBMaxOpenConns)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
