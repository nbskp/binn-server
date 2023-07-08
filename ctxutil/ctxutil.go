package ctxutil

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
)

var keyID = struct{}{}

func SetID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyID, id)
}

func ID(ctx context.Context) string {
	if id, ok := ctx.Value(keyID).(string); ok {
		fmt.Println(id)
		return id
	} else {
		fmt.Println("!!!!")
		return ""
	}
}

var keyLogAttrs = struct{}{}

func AddLogAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	if attrs, ok := ctx.Value(keyLogAttrs).([]slog.Attr); ok {
		attrs = append(attrs, newAttrs...)
		return withLogAttrs(ctx, attrs)
	} else {
		return withLogAttrs(ctx, newAttrs)
	}
}

func withLogAttrs(ctx context.Context, attrs []slog.Attr) context.Context {
	return context.WithValue(ctx, keyLogAttrs, attrs)
}

func LogAttrs(ctx context.Context) []slog.Attr {
	if attrs, ok := ctx.Value(keyLogAttrs).([]slog.Attr); ok {
		return attrs
	} else {
		return nil
	}
}
