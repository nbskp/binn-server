package auth

import (
	"context"
	"crypto/rand"
	"fmt"
)

type Provider interface {
	Issue(ctx context.Context, subID string) (token string, err error)
	Authorize(ctx context.Context, token string) (subID string, ok bool, err error)
}

func generateRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate a random string")
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
