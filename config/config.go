package config

import (
	"os"
	"strconv"
	"time"
)

const (
	envSendInterval = "BINN_SEND_INTERVAL_SEC"
	envRunEmitLoop  = "BINN_RUN_EMIT_LOOP"
)

var (
	defaultSendInterval = 10 * time.Second
	defaultRunEmitLoop  = true
)

type Config struct {
	RunEmitLoop  bool
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

	if v := os.Getenv(envRunEmitLoop); v == "false" {
		c.RunEmitLoop = false
	} else {
		c.RunEmitLoop = defaultRunEmitLoop
	}
	return c
}
