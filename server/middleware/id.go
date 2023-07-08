package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nbskp/binn-server/ctxutil"
	"github.com/nbskp/binn-server/logutil"
	"golang.org/x/exp/slog"
)

func IDMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uid, err := uuid.NewRandom(); err != nil {
			logger.Warn("failed to generate uuid")
		} else {
			ctx := ctxutil.SetID(r.Context(), uid.String())
			r = r.WithContext(ctxutil.AddLogAttrs(ctx, logutil.AttrID(uid.String())))
		}
		next.ServeHTTP(w, r)
	})
}
