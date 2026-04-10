package server

import (
	"context"
	"database/sql"
	"fmt"
	"mini-ledger/config"
	"mini-ledger/internal/middleware"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type server struct {
	httpServer *http.Server
	router     *chi.Mux
	huma       huma.API
	db         *sql.DB
	config     config.Config
}

func New(conf config.Config, db *sql.DB) *server {
	router := chi.NewRouter()

	srv := &http.Server{ // nolint:exhaustruct
		Addr:         fmt.Sprintf(":%d", conf.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(conf.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(conf.IdleTimeoutSec) * time.Second,
	}

	humaConfig := huma.DefaultConfig(conf.Name, conf.Version)

	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		if status >= 500 {
			log.Error().
				Int("status", status).
				Str("msg", msg).
				Errs("errors", errs).
				Msg("Server error occurred")
			return &huma.ErrorModel{ // nolint:exhaustruct
				Status: status,
				Title:  http.StatusText(status),
				Detail: http.StatusText(status),
			}
		}
		details := make([]*huma.ErrorDetail, len(errs))
		for i := range errs {
			if converted, ok := errs[i].(huma.ErrorDetailer); ok {
				details[i] = converted.ErrorDetail()
			}
		}
		return &huma.ErrorModel{ // nolint:exhaustruct
			Status: status,
			Title:  http.StatusText(status),
			Detail: msg,
			Errors: details,
		}
	}

	return &server{
		httpServer: srv,
		router:     router,
		huma:       humachi.New(router, humaConfig),
		db:         db,
		config:     conf,
	}
}

func (s *server) RegisterRoutes() {
	s.huma.UseMiddleware(middleware.Logger(s.config, s.huma))
	s.huma.UseMiddleware(middleware.Recoverer(s.config, s.huma))
	registerAllRoutes(s)
}

func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) GracefulShutdown() error {
	defer func() {
		err := s.db.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(s.config.GracefulShutdownTimeoutSec)*time.Second,
	)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
