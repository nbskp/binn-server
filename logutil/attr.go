package logutil

import (
	"github.com/nbskp/binn-server/binn"
	"golang.org/x/exp/slog"
)

const (
	attrIDKey = "id"

	attrBottleKey          = "bottle"
	attrBottleIDKey        = "id"
	attrBottleMsgKey       = "msg"
	attrBottleTokenKey     = "token"
	attrBottleExpiredAtKey = "expired_at"

	attrHTTPKey           = "http"
	attrHTTPMethodKey     = "method"
	attrHTTPPathKey       = "path"
	attrHTTPStatusCodeKey = "status_code"
)

func AttrID(id string) slog.Attr {
	return slog.String(attrIDKey, id)
}

func AttrBottle(b *binn.Bottle) slog.Attr {
	return slog.Group(attrBottleKey,
		slog.String(attrBottleIDKey, b.ID),
		slog.String(attrBottleMsgKey, b.Msg),
		slog.String(attrBottleTokenKey, b.Msg),
		slog.Int64(attrBottleExpiredAtKey, b.ExpiredAt),
	)
}

func AttrHTTP(method string, path string, statusCode int) slog.Attr {
	return slog.Group(attrHTTPKey,
		slog.String(attrHTTPMethodKey, method),
		slog.String(attrHTTPPathKey, path),
		slog.Int(attrHTTPStatusCodeKey, statusCode),
	)
}
