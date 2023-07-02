package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFromEnv(t *testing.T) {
	type args struct {
		envIntervalSec string
	}
	cases := []struct {
		name     string
		args     args
		expected Config
	}{
		{
			name: "env is filled",
			args: args{
				envIntervalSec: "30",
			},
			expected: Config{
				SendInterval: 30 * time.Second,
			},
		},
		{
			name: "env is not filled",
			args: args{
				envIntervalSec: "",
			},
			expected: Config{
				SendInterval: 10 * time.Second,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("BINN_SEND_INTERVAL_SEC", c.args.envIntervalSec)
			cfg := NewFromEnv()
			assert.Equal(t, c.expected, cfg)
		})
	}
}
