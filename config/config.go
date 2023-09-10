package config

import (
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

const (
	envPort = "PORT"

	envSendInterval           = "BINN_SEND_INTERVAL_SEC"
	envBottleExpiration       = "BINN_BOTTLE_EXPIRATION_SEC"
	envSubscriptionExpiration = "BINN_SUBSCRIPTION_EXPIRATION_SEC"
	envNumBottles             = "BINN_NUM_BOTTLES"

	envAuthKey = "AUTH_KEY"

	envRedisAddr           = "REDIS_ADDR"
	envRedisPassword       = "REDIS_PASSWORD"
	envRedisUsername       = "REDIS_USERNAME"
	envRedisBottleDB       = "REDIS_BOTTLE_DB"
	envRedisSubscriptionDB = "REDIS_SUBSCRIPTION_DB"
)

var (
	defaultPort                   = "8080"
	defaultSendInterval           = 10 * time.Second
	defaultNumBottles             = 10
	defaultBottleExpiration       = 60 * 10 * time.Second
	defaultSubscriptionExpiration = 60 * 15 * time.Second
	defaultRedisBottleDB          = 0
	defaultRedisSubscriptionDB    = 1
)

type Config struct {
	Port string

	SendInterval           time.Duration
	BottleExpiration       time.Duration
	SubscriptionExpiration time.Duration
	NumBottles             int

	AuthKey string

	RedisAddr           string
	RedisUsername       string
	RedisPassword       string
	RedisBottleDB       int
	RedisSubscriptionDB int
}

func NewFromEnv(logger *slog.Logger) Config {
	return Config{
		Port: loadPort(logger),

		SendInterval:           loadSendInterval(logger),
		BottleExpiration:       loadBottleExpiration(logger),
		SubscriptionExpiration: loadSubscriptionExpiration(logger),
		NumBottles:             loadNumBottles(logger),

		AuthKey: loadAuthKey(logger),

		RedisAddr:           loadRedisAddr(logger),
		RedisUsername:       loadRedisUsername(logger),
		RedisPassword:       loadRedisPassword(logger),
		RedisBottleDB:       loadRedisBottleDB(logger),
		RedisSubscriptionDB: loadRedisSubscriptionDB(logger),
	}
}

func loadPort(logger *slog.Logger) string {
	s := os.Getenv(envPort)
	if s == "" {
		logger.Warn("cannot load port from env, use default")
		return defaultPort
	}
	return s
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

func loadNumBottles(logger *slog.Logger) int {
	i, err := strconv.Atoi(os.Getenv(envNumBottles))
	if err != nil {
		logger.Warn("cannot load num of bottles from env, use default")
		return defaultNumBottles
	}
	return i
}

func loadAuthKey(logger *slog.Logger) string {
	s := os.Getenv(envAuthKey)
	if s == "" {
		logger.Error("cannot load auth key")
		os.Exit(0)
	}
	return s
}

func loadRedisAddr(logger *slog.Logger) string {
	s := os.Getenv(envRedisAddr)
	if s == "" {
		logger.Error("cannot load redis address")
		os.Exit(0)
	}
	return s
}

func loadRedisUsername(logger *slog.Logger) string {
	return os.Getenv(envRedisUsername)
}

func loadRedisPassword(logger *slog.Logger) string {
	return os.Getenv(envRedisPassword)
}

func loadRedisBottleDB(logger *slog.Logger) int {
	i, err := strconv.Atoi(os.Getenv(envRedisBottleDB))
	if err != nil {
		logger.Warn("cannot load redis bottle db from env, use default")
		return defaultRedisBottleDB
	}
	return i
}

func loadRedisSubscriptionDB(logger *slog.Logger) int {
	i, err := strconv.Atoi(os.Getenv(envRedisSubscriptionDB))
	if err != nil {
		logger.Warn("cannot load redis subscription db from env, use default")
		return defaultRedisSubscriptionDB
	}
	return i
}
