package auth

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenRedisProvider struct {
	cli  *redis.Client
	tLen int
}

func NewTokenRedisProvider(cli *redis.Client, tLen int) *tokenRedisProvider {
	return &tokenRedisProvider{
		cli:  cli,
		tLen: tLen,
	}
}

func (p *tokenRedisProvider) Issue(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return "", errors.New("timeout issue token")
		default:
		}
		token, err := generateRandomStr(uint32(p.tLen))
		if err != nil {
			return "", err
		}
		_, err = p.cli.Get(ctx, token).Result()
		if err != nil {
			if err == redis.Nil {
				if _, err := p.cli.Set(ctx, token, 0, -1).Result(); err != nil {
					return "", err
				}
				return token, nil
			}
			return "", err
		}
	}
}

func (p *tokenRedisProvider) Authorize(ctx context.Context, token string) (string, bool, error) {
	_, err := p.cli.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return token, true, nil
}
