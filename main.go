package main

import (
	"context"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/config"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server"
	"golang.org/x/exp/slog"

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
			Network: "tcp",
			Addr:    c.RedisAddr,
			DB:      0,
		}),
		10,
		c.BottleExpiration,
	)
	if err != nil {
		log.Fatal(err)
	}
	sh := binn.NewSubscriptionsRedisHandler(
		redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    c.RedisAddr,
			DB:      1,
		}),
		c.SubscriptionExpiration,
	)

	bn := binn.NewBinn(c.SendInterval, bh, sh)

	provider := auth.NewTokenRedisProvider(
		redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    c.RedisAddr,
			DB:      2,
		}),
		10,
	)
	srv := server.New(bn, provider, ":8080", l)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
