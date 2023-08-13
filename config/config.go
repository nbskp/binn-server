package config

import (
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

const (
	envSendInterval           = "BINN_SEND_INTERVAL_SEC"
	envBottleExpiration       = "BINN_BOTTLE_EXPIRATION_SEC"
	envSubscriptionExpiration = "BINN_SUBSCRIPTION_EXPIRATION_SEC"
	envRedisAddr              = "BINN_REDIS_ADDR"
)

var (
	defaultSendInterval           = 10 * time.Second
	defaultBottleExpiration       = 60 * 10 * time.Second
	defaultSubscriptionExpiration = 60 * 15 * time.Second
)

type Config struct {
	SendInterval           time.Duration
	BottleExpiration       time.Duration
	SubscriptionExpiration time.Duration

	RedisAddr string
}

func NewFromEnv(logger *slog.Logger) Config {
	c := Config{}

	c.SendInterval = loadSendInterval(logger)
	c.BottleExpiration = loadBottleExpiration(logger)
	c.SubscriptionExpiration = loadSubscriptionExpiration(logger)
	c.RedisAddr = loadRedisAddr(logger)

	return c
}

func loadSendInterval(logger *slog.Logger) time.Duration {
	i, err := strconv.Atoi(os.Getenv(envSendInterval))
	if err != nil {
		logger.Warn("cannot load send interval from env, use default")
		return defaultSendInterval
	}
	return time.Duration(i) * time.Second
}

func loadBottleExpiration(logger *slog.Logger) time.Duration {
	i, err := strconv.Atoi(os.Getenv(envBottleExpiration))
	if err != nil {
		logger.Warn("cannot load bottle expiration from env, use default")
		return defaultBottleExpiration
	}
	return time.Duration(i) * time.Second
}

func loadSubscriptionExpiration(logger *slog.Logger) time.Duration {
	i, err := strconv.Atoi(os.Getenv(envSubscriptionExpiration))
	if err != nil {
		logger.Warn("cannot load subscription expiration from env, use default")
		return defaultSubscriptionExpiration
	}
	return time.Duration(i) * time.Second
}

func loadRedisAddr(logger *slog.Logger) string {
	s := os.Getenv(envRedisAddr)
	if s == "" {
		logger.Error("cannot load redis address")
		os.Exit(0)
	}
	return s
}
