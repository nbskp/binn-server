package logutil

import (
	"context"

	"github.com/nbskp/binn-server/ctxutil"
	"golang.org/x/exp/slog"
)

const (
	attrIDKey = "id"

	attrEventKey          = "event"
	attrEventConnected    = "connected"
	attrEventDisconnected = "disconnected"
)

func LogConnected(ctx context.Context, logger *slog.Logger, msg string) {
	logger.Info(msg,
		slog.String(attrIDKey, ctxutil.ID(ctx)),
		slog.String(attrEventKey, attrEventConnected),
	)
}

func LogDisconnected(ctx context.Context, logger *slog.Logger, msg string) {
	logger.Info(msg,
		slog.String(attrIDKey, ctxutil.ID(ctx)),
		slog.String(attrEventKey, attrEventDisconnected),
	)
}

func LogWithID(ctx context.Context, loggerFunc func(string, ...any), msg string, args ...any) {
	args = append([]any{slog.String(attrIDKey, ctxutil.ID(ctx))}, args...)
	loggerFunc(msg, args...)
}
