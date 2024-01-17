package binn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Bottle_IsExpired(t *testing.T) {
	type args struct {
		expiredAt time.Time
		now       time.Time
	}
	cs := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "now > expiration, return true",
			args: args{
				expiredAt: time.Date(2023, time.June, 10, 11, 59, 0, 0, time.UTC),
				now:       time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
		{
			name: "now = expiration, return false",
			args: args{
				expiredAt: time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC),
				now:       time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		{
			name: "now < expiration, return false",
			args: args{
				expiredAt: time.Date(2023, time.June, 10, 12, 1, 0, 0, time.UTC),
				now:       time.Date(2023, time.June, 10, 12, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
	}
	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			b := Bottle{ExpiredAt: c.args.expiredAt}
			assert.Equal(t, c.expected, b.IsExpired(c.args.now))
		})
	}
}
