package ctxutil

import (
	"context"

	"github.com/nbskp/binn-server/binn"
	"golang.org/x/exp/slog"
)

type keyID struct{}

func SetID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyID{}, id)
}

func ID(ctx context.Context) string {
	if id, ok := ctx.Value(keyID{}).(string); ok {
		return id
	} else {
		return ""
	}
}

type keyLogAttrs struct{}

func AddLogAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	if attrs, ok := ctx.Value(keyLogAttrs{}).([]slog.Attr); ok {
		attrs = append(attrs, newAttrs...)
		return withLogAttrs(ctx, attrs)
	} else {
		return withLogAttrs(ctx, newAttrs)
	}
}

func withLogAttrs(ctx context.Context, attrs []slog.Attr) context.Context {
	return context.WithValue(ctx, keyLogAttrs{}, attrs)
}

func LogAttrs(ctx context.Context) []slog.Attr {
	if attrs, ok := ctx.Value(keyLogAttrs{}).([]slog.Attr); ok {
		return attrs
	} else {
		return nil
	}
}

type keySubscriptionID struct{}

func SetSubscriptionID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keySubscriptionID{}, id)
}
func SubscriptionID(ctx context.Context) string {
	if id, ok := ctx.Value(keySubscriptionID{}).(string); ok {
		return id
	} else {
		return ""
	}
}

type keyBottle struct{}

func SetBottle(ctx context.Context, b *binn.Bottle) context.Context {
	return context.WithValue(ctx, keyBottle{}, b)
}

func Bottle(ctx context.Context) *binn.Bottle {
	if b, ok := ctx.Value(keyBottle{}).(*binn.Bottle); ok {
		return b
	} else {
		return nil
	}
}
