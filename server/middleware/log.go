package middleware

import (
	"net/http"

	"github.com/nbskp/binn-server/logutil"
	"golang.org/x/exp/slog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func AccessLogMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w_ := &loggingResponseWriter{
			w,
			http.StatusOK,
		}
		next.ServeHTTP(w_, r)

		logger.InfoCtx(r.Context(), "ok",
			logutil.AttrHTTP(r.Method, r.URL.String(), w_.statusCode),
		)
	})
}
