package middleware

import (
	"mini-ledger/config"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

func Logger(conf config.Config, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		newCtx := log.With().
			Caller().
			Str("remote_addr", ctx.RemoteAddr()).Str("proto", ctx.Version().Proto).
			Str("method", ctx.Method()).
			Str("path", ctx.URL().Path).
			Logger().
			WithContext(ctx.Context())
		ctx = huma.WithContext(ctx, newCtx)

		log.Ctx(ctx.Context()).Info().Caller().Msg("Request received")

		reqStart := time.Now()

		next(ctx)

		reqDuration := time.Since(reqStart)

		status := ctx.Status()
		evt := log.Ctx(ctx.Context()).Info() //nolint:zerologlint
		msg := "request completed"
		if status >= 500 {
			evt = log.Ctx(ctx.Context()).Error() //nolint:zerologlint
			msg = "request failed"
		}
		if status >= 400 && status < 500 {
			msg = "request completed with client error"
		}
		evt.Int("status_code", status).Dur("duration", reqDuration).Msg(msg)
	}
}
