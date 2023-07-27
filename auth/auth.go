package auth

import (
	"crypto/rand"
	"fmt"
)

type Provider interface {
	Issue() (string, error)
	Authorize(string) (string, bool, error)
}

type tokenProvider struct {
	tokens map[string]struct{}
	tLen   int
}

func NewTokenProvider(tLen int) *tokenProvider {
	return &tokenProvider{
		tokens: make(map[string]struct{}),
		tLen:   tLen,
	}
}

func (p *tokenProvider) Issue() (string, error) {
	var token string
	for _, ok := p.tokens[token]; ok || token == ""; {
		var err error
		token, err = generateRandomStr(uint32(p.tLen))
		if err != nil {
			return "", err
		}
	}
	p.tokens[token] = struct{}{}
	return token, nil
}

func (p *tokenProvider) Authorize(token string) (string, bool, error) {
	if _, ok := p.tokens[token]; ok {
		return token, true, nil
	}
	return "", false, nil
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
