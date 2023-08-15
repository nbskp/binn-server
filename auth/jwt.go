package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type jwtProvider struct {
	key jwk.Key
	exp time.Duration
}

func (p *jwtProvider) Issue(ctx context.Context, subID string) (string, error) {
	tok, err := jwt.NewBuilder().
		JwtID(subID).
		Expiration(time.Now().Add(p.exp).UTC()).
		Issuer("github.com/nbskp/binn-server").
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, p.key))
	if err != nil {
		return "", err
	}
	return string(signed), nil
}

type payload struct {
	JTI string `json:"jti"`
	Exp int64  `json:"exp"`
}

func (p *jwtProvider) Authorize(ctx context.Context, src string) (string, bool, error) {
	plByte, err := jws.Verify([]byte(src), jws.WithKey(jwa.HS256, p.key))
	if err != nil {
		return "", false, err
	}
	var pl payload
	if err := json.NewDecoder(bytes.NewReader(plByte)).Decode(&pl); err != nil {
		return "", false, err
	}
	if time.Now().UTC().Unix() > pl.Exp {
		return "", false, nil
	}
	return pl.JTI, true, nil
}

func NewJWTProvider(key jwk.Key, exp time.Duration) *jwtProvider {
	return &jwtProvider{
		key: key,
		exp: exp,
	}
}
