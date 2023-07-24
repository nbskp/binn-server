package binn

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_BottleQueue_Push(t *testing.T) {
	nowTime := time.Now()
	now = func() time.Time {
		return nowTime
	}
	defer func() {
		now = time.Now
	}()
	notExpiredUnixSec := nowTime.Add(24 * time.Hour).Unix()
	expiredUnixSec := nowTime.Add(-24 * time.Hour).Unix()

	type args struct {
		b    Bottle
		sbs  []statefulBottle
		size int
	}
	type expected struct {
		sbs     []statefulBottle
		wantErr assert.ErrorAssertionFunc
	}
	cases := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "OK",
			args: args{
				b: Bottle{ID: "0", Msg: "new msg"},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: notExpiredUnixSec},
						state:  stateUnavailable,
					},
				},
				size: 1,
			},
			expected: expected{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "new msg", ExpiredAt: 0},
						state:  stateAvailable,
					},
				},
				wantErr: assert.NoError,
			},
		},
		{
			name: "NG: not found bottle",
			args: args{
				b: Bottle{ID: "1", Msg: "new msg"},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: notExpiredUnixSec},
						state:  stateUnavailable,
					},
				},
				size: 1,
			},
			expected: expected{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: notExpiredUnixSec},
						state:  stateUnavailable,
					},
				},
				wantErr: assert.Error,
			},
		},
		{
			name: "NG: available bottle",
			args: args{
				b: Bottle{ID: "0", Msg: "new msg"},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: expiredUnixSec},
						state:  stateAvailable,
					},
				},
				size: 1,
			},
			expected: expected{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: expiredUnixSec},
						state:  stateAvailable,
					},
				},
				wantErr: assert.Error,
			},
		},
		{
			name: "NG: expired bottle",
			args: args{
				b: Bottle{ID: "0", Msg: "new msg"},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: expiredUnixSec},
						state:  stateUnavailable,
					},
				},
				size: 1,
			},
			expected: expected{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "0", Msg: "", ExpiredAt: 0},
						state:  stateAvailable,
					},
				},
				wantErr: assert.Error,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q := &bottleQueue{
				sbs:  c.args.sbs,
				size: c.args.size,
			}
			err := q.Push(&c.args.b)
			c.expected.wantErr(t, err)
			assert.Equal(t, c.expected.sbs, q.sbs)
		})
	}
}

func Test_BottleQueue_Pop(t *testing.T) {
	nowTime := time.Now()
	now = func() time.Time {
		return nowTime
	}
	defer func() {
		now = time.Now
	}()
	expiredUnixSec := nowTime.Add(-24 * time.Hour).Unix()

	type args struct {
		sbs        []statefulBottle
		size       int
		expiration time.Duration
	}
	type expected struct {
		poppedBottles []*Bottle
		wantErrs      []assert.ErrorAssertionFunc
		sbs           []statefulBottle
	}
	cases := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "all bottles are available",
			args: args{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: 0},
						state:  stateAvailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: 0},
						state:  stateAvailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: 0},
						state:  stateAvailable,
					},
				},
				size:       3,
				expiration: 1 * time.Hour,
			},
			expected: expected{
				poppedBottles: []*Bottle{
					{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
				},
				wantErrs: []assert.ErrorAssertionFunc{
					assert.NoError,
					assert.NoError,
					assert.NoError,
				},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
				},
			},
		},
		{
			name: "has unavailable bottles",
			args: args{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: 0},
						state:  stateAvailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: 0},
						state:  stateAvailable,
					},
				},
				size:       3,
				expiration: 1 * time.Hour,
			},
			expected: expected{
				poppedBottles: []*Bottle{
					{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					nil,
				},
				wantErrs: []assert.ErrorAssertionFunc{
					assert.NoError,
					assert.NoError,
					assert.Error,
				},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
				},
			},
		},
		{
			name: "has expired bottles",
			args: args{
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: 0},
						state:  stateAvailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: expiredUnixSec},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: 0},
						state:  stateAvailable,
					},
				},
				size:       3,
				expiration: 1 * time.Hour,
			},
			expected: expected{
				poppedBottles: []*Bottle{
					{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
					{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
				},
				wantErrs: []assert.ErrorAssertionFunc{
					assert.NoError,
					assert.NoError,
					assert.NoError,
				},
				sbs: []statefulBottle{
					{
						bottle: &Bottle{ID: "1", Msg: "msg1", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "2", Msg: "msg2", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
					{
						bottle: &Bottle{ID: "3", Msg: "msg3", ExpiredAt: nowTime.Add(1 * time.Hour).Unix()},
						state:  stateUnavailable,
					},
				},
			},
		},
	}
	for _, c := range cases {
		fmt.Println()
		t.Run(c.name, func(t *testing.T) {
			q := &bottleQueue{
				sbs:        c.args.sbs,
				size:       c.args.size,
				expiration: c.args.expiration,
			}

			for i := 0; i < c.args.size; i++ {
				b, err := q.Pop()
				assert.Equal(t, c.expected.poppedBottles[i], b)
				c.expected.wantErrs[i](t, err)
			}
			assert.Equal(t, c.expected.sbs, q.sbs)
		})
	}
}
