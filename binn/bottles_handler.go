package binn

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type bottlesRedisHandler struct {
	cli  *redis.Client
	exp  time.Duration
	size int
}

func bottleID(n int) string {
	return strconv.Itoa(n)
}

func bottleKey(id string) string {
	return fmt.Sprintf("bottle:%s", id)
}

func bottleShadowKey(id string) string {
	return fmt.Sprintf("bottle:%s.shadow", id)
}

func (h *bottlesRedisHandler) Set(ctx context.Context, b *Bottle) error {
	_, err := h.cli.HGetAll(ctx, bottleShadowKey(b.ID)).Result()
	if err != nil {
		if err == redis.Nil {
			return NewBinnError(CodeNotFoundBottle, "not found the bottle", err)
		}
		return err
	}
	if _, err := h.cli.Del(ctx, bottleShadowKey(b.ID)).Result(); err != nil {
		return err
	}
	if _, err := h.cli.HSet(ctx, bottleKey(b.ID), "msg", b.Msg).Result(); err != nil {
		return err
	}
	return nil
}

func (h *bottlesRedisHandler) Next(ctx context.Context) (*Bottle, error) {
	var b *Bottle
	for i := 0; i < h.size; i++ {
		id := bottleID(i)
		ex, err := h.cli.Exists(ctx, bottleShadowKey(id)).Result()
		if err != nil {
			return nil, err
		}
		if ex == 0 {
			bv, err := h.cli.HGetAll(ctx, bottleKey(id)).Result()
			if err != nil {
				return nil, err
			}
			b = &Bottle{
				ID:  id,
				Msg: bv["msg"],
			}
			break
		}
	}
	if b == nil {
		return nil, nil
	}
	sKey := bottleShadowKey(b.ID)
	if _, err := h.cli.HSet(ctx, sKey, 0, 0).Result(); err != nil {
		return nil, err
	}
	expiredAt := now().Add(h.exp)
	if _, err := h.cli.ExpireAt(ctx, sKey, expiredAt).Result(); err != nil {
		return nil, err
	}
	b.ExpiredAt = expiredAt
	return b, nil
}

func NewBottlesRedisHandler(ctx context.Context, cli *redis.Client, size int, exp time.Duration) (*bottlesRedisHandler, error) {
	ks, err := cli.Keys(ctx, "bottle:*").Result()
	if err != nil {
		return nil, err
	}
	if len(ks) == 0 {
		for i := 0; i < size; i++ {
			_, err := cli.HSet(ctx, bottleKey(bottleID(i)), "msg", "").Result()
			if err != nil {
				return nil, err
			}
		}
	}
	return &bottlesRedisHandler{
		cli:  cli,
		exp:  exp,
		size: size,
	}, nil
}
