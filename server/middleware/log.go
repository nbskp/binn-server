package middleware

import (
	"net/http"

	"github.com/nbskp/binn-server/logutil"
	"golang.org/x/exp/slog"
)

func LogConnectionEventMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.InfoCtx(r.Context(), "connected", logutil.AttrEventConnected())
		}
		next.ServeHTTP(w, r)
		if logger != nil {
			logger.InfoCtx(r.Context(), "disconnected", logutil.AttrEventDisconnected())
		}
	})
}
