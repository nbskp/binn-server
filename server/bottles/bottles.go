package bottles

import (
	"encoding/json"
	"net/http"

	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/ctxutil"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server/middleware"
	"golang.org/x/exp/slog"
)

func NewBottlesMux(bn *binn.Binn, auth auth.Provider, logger *slog.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/", middleware.AuthMiddleware(http.HandlerFunc(bottlesHandlerFunc(bn, logger)), auth, logger))
	return r
}

func bottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	postHf := postBottlesHandlerFunc(bn, logger)
	getHf := getBottlesHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHf(w, r)
		case http.MethodPost:
			postHf(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID := ctxutil.SubscriptionID(r.Context())
		b, err := bn.GetBottle(subID)
		if err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if b == nil {
			logger.InfoCtx(r.Context(), "no bottles")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		resp := ToResponse(b)
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.InfoCtx(r.Context(), "send bottle", logutil.AttrEventSendBottle(), logutil.AttrBottle(b))
	}
}

func postBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		b := reqBody.ToBottle()
		subID := ctxutil.SubscriptionID(r.Context())
		if err := bn.Publish(subID, b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		logger.InfoCtx(r.Context(), "receive bottle", logutil.AttrEventReceiveBottle(), logutil.AttrBottle(b))
		w.WriteHeader(http.StatusOK)
	}
}
