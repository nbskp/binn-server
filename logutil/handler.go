package logutil

import (
	"context"

	"github.com/nbskp/binn-server/ctxutil"
	"golang.org/x/exp/slog"
)

type ctxHandler struct {
	baseHandler slog.Handler
}

func (h *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctxHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
	}
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	return &ctxHandler{
		baseHandler: h.baseHandler.WithGroup(name),
	}
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := ctxutil.LogAttrs(ctx)
	return h.baseHandler.WithAttrs(attrs).Handle(ctx, r)
}

func NewCtxHandler(h slog.Handler) *ctxHandler {
	return &ctxHandler{
		baseHandler: h,
	}
}
