package middleware

import (
	"net/http"
	"strings"

	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/ctxutil"
	"golang.org/x/exp/slog"
)

func AuthMiddleware(next http.Handler, provider auth.Provider, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		es := strings.Split(r.Header.Get("Authorization"), " ")
		if len(es) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		schema, token := es[0], es[1]
		if schema != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		subID, ok, err := provider.Authorize(r.Context(), token)
		if err != nil || !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(ctxutil.SetSubscriptionID(r.Context(), subID))
		next.ServeHTTP(w, r)
	})
}
