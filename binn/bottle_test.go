package binn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Bottle_IsExpired(t *testing.T) {
	now = func() time.Time {
		return time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC)
	}
	defer func() {
		now = time.Now
	}()
	cs := []struct {
		name     string
		args     Bottle
		expected bool
	}{
		{
			name: "now > expiration, return true",
			args: Bottle{
				ExpiredAt: time.Date(2023, time.June, 10, 11, 59, 0, 0, time.UTC),
			},
			expected: true,
		},
		{
			name: "now = expiration, return false",
			args: Bottle{
				ExpiredAt: time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		{
			name: "now < expiration, return false",
			args: Bottle{
				ExpiredAt: time.Date(2023, time.June, 10, 12, 1, 0, 0, time.UTC),
			},
			expected: false,
		},
	}
	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.args.IsExpired())
		})
	}
}
