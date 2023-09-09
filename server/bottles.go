package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/ctxutil"
	"github.com/nbskp/binn-server/server/middleware"
	"golang.org/x/exp/slog"
)

func NewBottlesMux(bn *binn.Binn, provider auth.Provider, logger *slog.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/", middleware.AuthMiddleware(http.HandlerFunc(bottlesHandlerFunc(bn, logger)), provider, logger))
	return r
}

func subscribeBottlesHandler(bn *binn.Binn, provider auth.Provider, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		subID := uuid.New().String()
		token, err := provider.Issue(r.Context(), subID)
		if err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bn.Subscribe(r.Context(), subID)
		if err := json.NewEncoder(w).Encode(newSubscribeBottlesResponse(token)); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
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
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID := ctxutil.SubscriptionID(r.Context())
		b, err := bn.GetBottle(r.Context(), subID)
		if err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		if b == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		resp := toBottlesResponse(b)
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func postBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody bottlesRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqBody); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b := reqBody.toBottles()
		subID := ctxutil.SubscriptionID(r.Context())
		if err := bn.Publish(r.Context(), subID, b); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
