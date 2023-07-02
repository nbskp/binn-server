package server

import (
	"log"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles"
	"github.com/nbskp/binn-server/server/ping"
	"github.com/nbskp/binn-server/server/stream"
)

func New(bn *binn.Binn, addr string, logger *log.Logger) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("/bottles/stream", stream.HandlerFunc(bn, logger))
	r.HandleFunc("/bottles", bottles.HandlerFunc(bn, logger))
	r.HandleFunc("/ping", ping.HandlerFunc())

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
