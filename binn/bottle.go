package binn

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Bottle struct {
	ID        string
	Msg       string
	Token     string
	ExpiredAt time.Time
}

func (b *Bottle) IsExpired(now time.Time) bool {
	return b.ExpiredAt.Before(now)
}

func generateToken(digit uint32) (string, error) {
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
