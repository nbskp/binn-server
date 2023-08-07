package auth

import (
	"context"
	"crypto/rand"
	"fmt"
)

type Provider interface {
	Issue(context.Context) (string, error)
	Authorize(context.Context, string) (string, bool, error)
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
