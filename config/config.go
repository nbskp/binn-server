package config

import (
	"os"
	"strconv"
	"time"
)

const (
	envSendInterval = "BINN_SEND_INTERVAL_SEC"
)

var (
	defaultSendInterval = 10 * time.Second
)

type Config struct {
	SendInterval time.Duration
}

func NewFromEnv() Config {
	c := Config{}
	i, err := strconv.Atoi(os.Getenv(envSendInterval))
	if err != nil {
		c.SendInterval = defaultSendInterval
	} else {
		c.SendInterval = time.Duration(i) * time.Second
	}
	return c
}
