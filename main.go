package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/config"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server"
	"golang.org/x/exp/slog"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/redis/go-redis/v9"
)

var programLevel = new(slog.LevelVar)

func main() {
	l := slog.New(logutil.NewCtxHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))
	c := config.NewFromEnv(l)

	ctx := context.Background()
	bh, err := binn.NewBottlesRedisHandler(
		ctx,
		redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     c.RedisAddr,
			Username: c.RedisUsername,
			Password: c.RedisPassword,
			DB:       c.RedisBottleDB,
		}),
		c.NumBottles,
		c.BottleExpiration,
	)
	if err != nil {
		l.Error(fmt.Sprintf("initializing bottles handler is failed: %v", err))
		os.Exit(0)
	}
	sh := binn.NewSubscriptionsRedisHandler(
		redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     c.RedisAddr,
			Username: c.RedisUsername,
			Password: c.RedisPassword,
			DB:       c.RedisSubscriptionDB,
		}),
		c.SubscriptionExpiration,
	)

	bn := binn.NewBinn(c.SendInterval, bh, sh)

	key, err := jwk.FromRaw([]byte(c.AuthKey))
	if err != nil {
		l.Error(fmt.Sprintf("initializing auth key is failed: %v", err))
		os.Exit(0)
	}
	provider := auth.NewJWTProvider(key, c.SubscriptionExpiration)
	srv := server.New(bn, provider, fmt.Sprintf(":%s", c.Port), l)
	if err := srv.ListenAndServe(); err != nil {
		l.Error(fmt.Sprintf("running server is failed: %v", err))
		os.Exit(0)
	}
}
