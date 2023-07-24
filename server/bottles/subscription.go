package bottles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles/response"
	"golang.org/x/exp/slog"
)

func newSubscriptionsMux(bn *binn.Binn, logger *slog.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/", subscriptionsHandler(bn, logger))
	return r
}

func subscriptionsHandler(bn *binn.Binn, logger *slog.Logger) http.Handler {
	subs := map[string][]*binn.Bottle{}
	postHf := postSubscription(bn, logger, subs)
	getHf := getSubscription(bn, logger, subs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHf.ServeHTTP(w, r)
		case http.MethodGet:
			getHf.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

type subscriptionRequest struct {
	ID string `json:"id"`
}

func postSubscription(bn *binn.Binn, logger *slog.Logger, subs map[string][]*binn.Bottle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var subReq subscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		subs[subReq.ID] = []*binn.Bottle{}
		bn.Subscribe(func(b *binn.Bottle) bool {
			fmt.Println("emited")
			subs[subReq.ID] = append(subs[subReq.ID], b)
			return true
		})
		logger.InfoCtx(r.Context(), "subscribed")
		return
	})
}

type subscriptionResponse struct {
	Bottles []*response.Response `json:"bottles"`
}

func getSubscription(bn *binn.Binn, logger *slog.Logger, subs map[string][]*binn.Bottle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if strings.IndexRune(id, '/') != -1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		bs, ok := subs[id]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := subscriptionResponse{Bottles: make([]*response.Response, len(bs))}
		for i, b := range bs {
			res.Bottles[i] = response.ToResponse(b)
		}
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.ErrorCtx(r.Context(), err.Error())
			return
		}
		subs[id] = []*binn.Bottle{}
	})
}
