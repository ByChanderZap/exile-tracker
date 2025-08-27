package utils

import (
	"net/http"

	"github.com/rs/zerolog"
)

// ZerologMiddleware returns a middleware that logs HTTP requests using the provided zerolog.Logger
func ZerologMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	childLog := ChildLogger("http").With().Logger()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			childLog.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Msg("HTTP request")
			next.ServeHTTP(w, r)
		})
	}
}
