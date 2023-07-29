package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFromEnv(t *testing.T) {
	type args struct {
		envIntervalSec               string
		envSubscriptionExpirationSec string
	}
	cases := []struct {
		name     string
		args     args
		expected Config
	}{
		{
			name: "env is filled",
			args: args{
				envIntervalSec:               "30",
				envSubscriptionExpirationSec: "10",
			},
			expected: Config{
				SendInterval:           30 * time.Second,
				SubscriptionExpiration: 10 * time.Second,
			},
		},
		{
			name: "env is not filled",
			args: args{
				envIntervalSec:               "",
				envSubscriptionExpirationSec: "",
			},
			expected: Config{
				SendInterval:           10 * time.Second,
				SubscriptionExpiration: 60 * 15 * time.Second,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("BINN_SEND_INTERVAL_SEC", c.args.envIntervalSec)
			t.Setenv("BINN_SUBSCRIPTION_EXPIRATION_SEC", c.args.envSubscriptionExpirationSec)
			cfg := NewFromEnv()
			assert.Equal(t, c.expected, cfg)
		})
	}
}
