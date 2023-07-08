package logutil

import (
	"context"

	"github.com/nbskp/binn-server/ctxutil"
	"golang.org/x/exp/slog"
)

const (
	attrIDKey = "id"

	attrEventKey               = "event"
	attrEventConnectedValue    = "connected"
	attrEventDisconnectedValue = "disconnected"
)

func AttrEventConnected() slog.Attr {
	return slog.String(attrEventKey, attrEventConnectedValue)
}

func AttrEventDisconnected() slog.Attr {
	return slog.String(attrEventKey, attrEventDisconnectedValue)
}

func AttrID(id string) slog.Attr {
	return slog.String(attrIDKey, id)
}

func SetID(ctx context.Context, id string) context.Context {
	return ctxutil.AddLogAttrs(ctx, AttrID(id))
}
