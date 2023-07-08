package bottles

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/logutil"
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
			logger.ErrorCtx(r.Context(), "failed to cast http.ResponseWriter to flusher")
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
				logger.ErrorCtx(r.Context(), fmt.Sprintf("failed to encode a bottle: %v", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				closed <- struct{}{}
				return false
			}
			flusher.Flush()
			logger.InfoCtx(r.Context(), "send a bottle", logutil.AttrBottle(b), logutil.AttrEventSendBottle())
			return true
		})

		select {
		case <-r.Context().Done():
		case <-closed:
		}
		return
	}
}
