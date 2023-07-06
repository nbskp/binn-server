package ctxutil

import (
	"context"
	"fmt"
)

var (
	keyID = struct{}{}
)

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
