package auth

import "context"

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

func (p *tokenProvider) Issue(ctx context.Context) (string, error) {
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

func (p *tokenProvider) Authorize(ctx context.Context, token string) (string, bool, error) {
	if _, ok := p.tokens[token]; ok {
		return token, true, nil
	}
	return "", false, nil
}
