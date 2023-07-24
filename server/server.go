package server

import (
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles"
	"github.com/nbskp/binn-server/server/middleware"
	"github.com/nbskp/binn-server/server/ping"
	"golang.org/x/exp/slog"
)

func New(bn *binn.Binn, addr string, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: newHandler(bn, logger),
	}
}

func newHandler(bn *binn.Binn, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	r.Handle("/ping", http.HandlerFunc(ping.HandlerFunc()))
	r.Handle("/bottles/", http.StripPrefix("/bottles", bottles.NewBottlesMux(bn, logger)))
	r.Handle("/bottles:emit", emitHandler(bn, logger))
	return middleware.HTTPInfoMiddleware(middleware.IDMiddleware(r, logger))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST")
		next.ServeHTTP(w, r)
	})
}

func emitHandler(bn *binn.Binn, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := bn.Emit(); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
