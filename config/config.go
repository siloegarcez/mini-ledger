package config

import (
	"errors"
	"fmt"
)

type Config struct {
	Port                       int    `default:"80"                help:"Port to listen on."                                                                                       short:"p"`
	Version                    string `default:"development build" help:"Set the application version."                                                                             short:"v"`
	Name                       string `default:"Mini Ledger API"   help:"Set the application name."                                                                                short:"n"`
	Environment                string `default:"dev"               help:"Define whether the app is running in production."                                                         short:"e"`
	GracefulShutdownTimeoutSec int    `default:"5"                 help:"Define time in seconds for all goroutines to finish handling the current request during server shutdown." short:"t"`
	ReadTimeoutSec             int    `default:"15"                help:"Define time in seconds for request read timeout."                                                         short:"r"`
	WriteTimeoutSec            int    `default:"15"                help:"Define time in seconds for response write timeout."                                                       short:"w"`
	IdleTimeoutSec             int    `default:"60"                help:"Define time in seconds for idle connection timeout."                                                      short:"i"`
	DBUsername                 string `default:"dev"               help:"Set the database username."`
	DBPassword                 string `default:"dev"               help:"Set the database password."`
	DBHost                     string `default:"localhost"         help:"Set the database host."`
	DBPort                     string `default:"5432"              help:"Set the database port."`
	DBName                     string `default:"dev"               help:"Set the database name."`
	DBSchema                   string `default:"public"            help:"Set the PostgreSQL database schema."`
	DBConnMaxLifetimeSec       int    `default:"0"                 help:"Define maximum lifetime in seconds for a database connection."`
	DBMaxIdleConns             int    `default:"50"                help:"Define maximum number of idle database connections."`
	DBMaxOpenConns             int    `default:"50"                help:"Define maximum number of open database connections."`
	MigrationsPath             string `default:"./migrations"      help:"Set the database migrations directory."`
	LogLevel                   int    `default:"0"                 help:"Define log level (debug=0, info=1, warn=2, error=3, fatal=4, panic=5, nolevel=6, disabled=7)."`
	StartupTimeoutSec          int    `default:"10"                help:"Define time in seconds for server startup timeout."`
}

var (
	ErrInvalidConfig = errors.New("invalid configuration")
)

func (c Config) Validate() []error {
	errs := []error{}
	switch c.Environment {
	case EnvDev, EnvStaging, EnvProd:
		// ok
	default:
		errs = append(
			errs,
			fmt.Errorf(
				"%w: invalid environment '%s', must be one of: %s, %s, %s",
				ErrInvalidConfig,
				c.Environment,
				EnvDev,
				EnvStaging,
				EnvProd,
			),
		)
	}
	return errs
}
