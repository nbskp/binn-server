package middleware

import (
	"log"
	"net/http"
)

func LogConnectionEventMiddleware(next http.HandlerFunc, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Printf(`{"type":"connection","event":"connected"}`)
		}
		next(w, r)
		if logger != nil {
			logger.Printf(`{"type:"connection","event":"disconnected"}`)
		}
	})
}
