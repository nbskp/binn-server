package config

import (
	"os"
	"strconv"
	"time"
)

const (
	envSendInterval           = "BINN_SEND_INTERVAL_SEC"
	envSubscriptionExpiration = "BINN_SUBSCRIPTION_EXPIRATION_SEC"
)

var (
	defaultSendInterval           = 10 * time.Second
	defaultSubscriptionExpiration = 60 * 15 * time.Second
)

type Config struct {
	SendInterval           time.Duration
	SubscriptionExpiration time.Duration
}

func NewFromEnv() Config {
	c := Config{}

	i, err := strconv.Atoi(os.Getenv(envSendInterval))
	if err != nil {
		c.SendInterval = defaultSendInterval
	} else {
		c.SendInterval = time.Duration(i) * time.Second
	}

	se, err := strconv.Atoi(os.Getenv(envSubscriptionExpiration))
	if err != nil {
		c.SubscriptionExpiration = defaultSubscriptionExpiration
	} else {
		c.SubscriptionExpiration = time.Duration(se) * time.Second
	}

	return c
}
