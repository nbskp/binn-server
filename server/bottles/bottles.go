package bottles

import (
	"encoding/json"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server/bottles/request"
	"github.com/nbskp/binn-server/server/middleware"
	"golang.org/x/exp/slog"
)

func NewBottlesMux(bn *binn.Binn, logger *slog.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/", http.HandlerFunc(bottlesHandlerFunc(bn, logger)))
	r.Handle("/stream", middleware.LogConnectionEventMiddleware(StreamHandlerFunc(bn, logger), logger))
	r.Handle("/ws", middleware.LogConnectionEventMiddleware(WebsocketHandlerFunc(bn, logger), logger))

	r.Handle("/subscriptions/", http.StripPrefix("/subscriptions", newSubscriptionsMux(bn, logger)))
	return r
}

func bottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	postHf := postBottlesHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHf(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody request.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		b := reqBody.ToBottle()
		if err := bn.Publish(b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		logger.InfoCtx(r.Context(), "receive bottle", logutil.AttrEventReceiveBottle(), logutil.AttrBottle(b))
		w.WriteHeader(http.StatusOK)
	}
}
