package server

import (
	"encoding/json"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"golang.org/x/exp/slog"
)

func NewBottlesMux(bn *binn.Binn, logger *slog.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/", http.HandlerFunc(bottlesHandlerFunc(bn, logger)))
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
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getBottlesHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := bn.Get(r.Context())
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
		if err := bn.Set(r.Context(), b); err != nil {
			logger.ErrorCtx(r.Context(), err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
