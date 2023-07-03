package bottles

import (
	"encoding/json"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles/response"
	"golang.org/x/exp/slog"
)

func StreamHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	hf := getStreamHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			hf(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getStreamHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		flusher, ok := w.(http.Flusher)
		if !ok {
			logger.Error("failed to cast to flusher")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		closed := make(chan struct{}, 0)
		bn.Subscribe(func(b *binn.Bottle) bool {
			select {
			case <-r.Context().Done():
				return false
			default:
			}

			if err := enc.Encode(response.ToResponse(b)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				closed <- struct{}{}
				return false
			}
			flusher.Flush()
			return true
		})

		select {
		case <-r.Context().Done():
		case <-closed:
		}
		return
	}
}
