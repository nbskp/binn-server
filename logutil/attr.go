package logutil

import (
	"github.com/nbskp/binn-server/binn"
	"golang.org/x/exp/slog"
)

const (
	attrIDKey = "id"

	attrEventKey               = "event"
	attrEventConnectedValue    = "connected"
	attrEventDisconnectedValue = "disconnected"
	attrEventSendBottle        = "send-bottle"
	attrEventReceiveBottle     = "receive-bottle"

	attrBottleKey          = "bottle"
	attrBottleIDKey        = "id"
	attrBottleMsgKey       = "msg"
	attrBottleTokenKey     = "token"
	attrBottleExpiredAtKey = "expired_at"

	attrHTTPKey       = "http"
	attrHTTPMethodKey = "method"
	attrHTTPPathKey   = "path"
)

func AttrEventConnected() slog.Attr {
	return slog.String(attrEventKey, attrEventConnectedValue)
}

func AttrEventDisconnected() slog.Attr {
	return slog.String(attrEventKey, attrEventDisconnectedValue)
}

func AttrEventSendBottle() slog.Attr {
	return slog.String(attrEventKey, attrEventSendBottle)
}

func AttrEventReceiveBottle() slog.Attr {
	return slog.String(attrEventKey, attrEventReceiveBottle)
}

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

func AttrHTTP(method string, path string) slog.Attr {
	return slog.Group(attrHTTPKey,
		slog.String(attrHTTPMethodKey, method),
		slog.String(attrHTTPPathKey, path),
	)
}
