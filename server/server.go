package server

import (
	"net/http"

	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/middleware"
	"github.com/nbskp/binn-server/server/ping"
	"golang.org/x/exp/slog"
)

func New(bn *binn.Binn, auth auth.Provider, addr string, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: newHandler(bn, auth, logger),
	}
}

func newHandler(bn *binn.Binn, auth auth.Provider, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	r.Handle("/ping", http.HandlerFunc(ping.HandlerFunc()))

	r.Handle("/bottles", NewBottlesMux(bn, auth, logger))
	r.Handle("/bottles:subscribe", subscribeBottlesHandler(bn, auth, logger))

	return middleware.AccessLogMiddleware(
		middleware.IDMiddleware(middleware.CORSMiddleware(r), logger), logger)
}
