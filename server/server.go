package server

import (
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles"
	"github.com/nbskp/binn-server/server/ping"
	"golang.org/x/exp/slog"
)

func New(bn *binn.Binn, addr string, logger *slog.Logger) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("/ping", ping.HandlerFunc())
	r.Handle("/bottles/", http.StripPrefix("/bottles", bottles.NewBottlesMux(bn, logger)))

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST")
		next.ServeHTTP(w, r)
	})
}
