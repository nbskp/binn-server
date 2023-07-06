package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nbskp/binn-server/binn"
)

type response struct {
	ID        string `json:"id"`
	Msg       string `json:"message"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

func toResponse(b *binn.Bottle) *response {
	return &response{
		ID:        b.ID,
		Msg:       b.Msg,
		Token:     b.Token,
		ExpiredAt: b.ExpiredAt,
	}
}

func StreamHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	hf := getStreamHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			hf(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getStreamHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		flusher, ok := w.(http.Flusher)
		if !ok {
			logger.Printf("failed to cast to flusher")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if logger != nil {
			logger.Printf("[out] connected")
		}
		enc := json.NewEncoder(w)
		closed := make(chan struct{}, 0)
		bn.Subscribe(func(b *binn.Bottle) bool {
			select {
			case <-r.Context().Done():
				return false
			default:
			}

			if err := enc.Encode(toResponse(b)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				closed <- struct{}{}
				return false
			}
			flusher.Flush()
			if logger != nil {
				logger.Printf("[out] {\"message\": \"%s\"}\n", b.Msg)
			}
			flusher.Flush()
			return true
		})

		select {
		case <-r.Context().Done():
		case <-closed:
		}
		if logger != nil {
			logger.Printf("[out] disconnected")
		}
		return
	}
}
