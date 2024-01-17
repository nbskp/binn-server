package server

import (
	"net/http"

	"github.com/nbskp/binn-server/binn"
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
	r.Handle("/bottles", NewBottlesMux(bn, logger))

	return middleware.AccessLogMiddleware(
		middleware.IDMiddleware(middleware.CORSMiddleware(r), logger), logger)
}
