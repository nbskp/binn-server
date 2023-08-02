package binn

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stubQueue struct{}

func (s *stubQueue) Set(b *Bottle) error {
	return nil
}

func (s *stubQueue) Next() (*Bottle, error) {
	return &Bottle{ID: "1", Msg: "sample"}, nil
}

func Test_Binn_GetBottle(t *testing.T) {
	nowTime := time.Now()
	now = func() time.Time {
		return nowTime
	}
	defer func() { now = time.Now }()

	q := &stubQueue{}
	itv := time.Duration(1)
	subExp := time.Duration(10)
	bn := &Binn{
		bh: q,
		subs: map[string]*Subscription{
			"example_id": {
				id:        "example_id",
				expiredAt: now().Add(subExp),
				nextTime:  now().Add(itv),
				bottleIDs: map[string]struct{}{},
			},
		},
		itv: itv,
	}

	b, err := bn.GetBottle("example_id")
	assert.NoError(t, err)
	assert.Nil(t, b)

	now = func() time.Time {
		return nowTime.Add(time.Duration(1))
	}
	fmt.Println(now())

	b, err = bn.GetBottle("example_id")
	assert.NoError(t, err)
	assert.Equal(t, &Bottle{ID: "1", Msg: "sample"}, b)
}

func Test_binn_Publish(t *testing.T) {
	nowTime := time.Now()
	now = func() time.Time {
		return nowTime
	}
	defer func() { now = time.Now }()

	q := &stubQueue{}
	itv := time.Duration(1)
	subExp := time.Duration(10)
	bn := &Binn{
		bh: q,
		subs: map[string]*Subscription{
			"example_id": {
				id:        "example_id",
				expiredAt: now().Add(subExp),
				nextTime:  now().Add(itv),
				bottleIDs: map[string]struct{}{"1": struct{}{}},
			},
		},
		itv: itv,
	}

	err := bn.Publish("example_id", &Bottle{ID: "1", Msg: "sample"})
	assert.NoError(t, err)

	err = bn.Publish("example_id", &Bottle{ID: "2", Msg: "sample"})
	assert.Error(t, err)

	err = bn.Publish("example_id", &Bottle{ID: "1", Msg: "sample"})
	assert.Error(t, err)
}
