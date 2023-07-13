package middleware

import (
	"net/http"

	"github.com/nbskp/binn-server/ctxutil"
	"github.com/nbskp/binn-server/logutil"
)

func HTTPInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctxutil.AddLogAttrs(r.Context(), logutil.AttrHTTP(r.Method, r.URL.Path)))
		next.ServeHTTP(w, r)
	})
}
