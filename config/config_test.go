package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFromEnv(t *testing.T) {
	type args struct {
		envRunEmitLoop string
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
				envRunEmitLoop: "false",
				envIntervalSec: "30",
			},
			expected: Config{
				RunEmitLoop:  false,
				SendInterval: 30 * time.Second,
			},
		},
		{
			name: "env is not filled",
			args: args{
				envRunEmitLoop: "",
				envIntervalSec: "",
			},
			expected: Config{
				RunEmitLoop:  true,
				SendInterval: 10 * time.Second,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("BINN_RUN_EMIT_LOOP", c.args.envRunEmitLoop)
			t.Setenv("BINN_SEND_INTERVAL_SEC", c.args.envIntervalSec)
			cfg := NewFromEnv()
			assert.Equal(t, c.expected, cfg)
		})
	}
}
