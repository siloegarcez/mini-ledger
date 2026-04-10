package database

import (
	"context"
	"database/sql"
	"fmt"
	"mini-ledger/config"
	"time"
)

const (
	sslDisabledConn = "postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s"
	sslEnabledConn  = "postgres://%s:%s@%s:%s/%s?sslmode=require&search_path=%s"
)

func Connect(ctx context.Context, conf config.Config) (*sql.DB, error) {
	connFormat := sslDisabledConn

	if conf.Environment == config.EnvProd || conf.Environment == config.EnvStaging {
		connFormat = sslEnabledConn
	}

	connStr := fmt.Sprintf(
		connFormat,
		conf.DBUsername,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
		conf.DBSchema,
	)

	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(conf.DBConnMaxLifetimeSec) * time.Second)
	db.SetMaxIdleConns(conf.DBMaxIdleConns)
	db.SetMaxOpenConns(conf.DBMaxOpenConns)

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
