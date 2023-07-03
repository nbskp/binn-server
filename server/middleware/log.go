package middleware

import (
	"net/http"

	"golang.org/x/exp/slog"
)

func LogConnectionEventMiddleware(next http.HandlerFunc, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Info("connected", "type", "connection", "event", "connected")
		}
		next(w, r)
		if logger != nil {
			logger.Info("disconnected", "type", "connection", "event", "disconnected")
		}
	})
}
