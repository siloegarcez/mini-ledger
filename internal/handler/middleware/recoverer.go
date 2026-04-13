package middleware

import (
	"mini-ledger/config"
	"net/http"
	"runtime/debug"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

func Recoverer(conf config.Config, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			if r := recover(); r != nil {
				log.Ctx(ctx.Context()).
					Error().
					Bytes("debug_stack", debug.Stack()).
					Interface("panic", r).
					Msg("Panic recovered")
				err := huma.WriteErr(
					api,
					ctx,
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
				)
				if err != nil {
					log.Ctx(ctx.Context()).Error().Err(err).Msg("Failed to write error response")
					ctx.SetStatus(http.StatusInternalServerError)
				}
			}
		}()
		next(ctx)
	}
}
