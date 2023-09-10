package binn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type subscriptionsRedisHandler struct {
	cli *redis.Client
	exp time.Duration
}

func subscriptionKey(id string) string {
	return fmt.Sprintf("subscription:%s", id)
}

func subscriptionToHashFields(s *Subscription) []interface{} {
	fs := make([]interface{}, 4)
	fs = append(fs, "next_time", s.nextTime)
	fs = append(fs, "bottle_ids", strings.Join(s.bottleIDs, ","))
	return fs
}

func mapToSubscription(m map[string]string) (*Subscription, error) {
	nt, err := time.Parse(time.RFC3339, m["next_time"])
	if err != nil {
		return nil, err
	}
	return &Subscription{
		nextTime:  nt,
		bottleIDs: strings.Split(m["bottle_ids"], ","),
	}, nil
}

func (sh *subscriptionsRedisHandler) Get(ctx context.Context, id string) (*Subscription, error) {
	vs, err := sh.cli.HGetAll(ctx, subscriptionKey(id)).Result()
	if err != nil {
		return nil, err
	}
	// https://github.com/redis/go-redis/issues/1668
	if len(vs) == 0 {
		return nil, nil
	}
	s, err := mapToSubscription(vs)
	if err != nil {
		return nil, err
	}
	s.id = id
	return s, nil
}

func (sh *subscriptionsRedisHandler) Update(ctx context.Context, sub *Subscription) error {
	fs := subscriptionToHashFields(sub)
	_, err := sh.cli.HSet(ctx, subscriptionKey(sub.id), fs...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (sh *subscriptionsRedisHandler) Add(ctx context.Context, sub *Subscription) error {
	fs := subscriptionToHashFields(sub)
	_, err := sh.cli.HSet(ctx, subscriptionKey(sub.id), fs...).Result()
	if err != nil {
		return err
	}
	_, err = sh.cli.ExpireAt(ctx, subscriptionKey(sub.id), now().Add(sh.exp)).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewSubscriptionsRedisHandler(cli *redis.Client, exp time.Duration) *subscriptionsRedisHandler {
	return &subscriptionsRedisHandler{cli: cli, exp: exp}
}
