package binn

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stubQueue struct{}

func (s *stubQueue) Pop() (*Bottle, error) {
	return &Bottle{ID: "1", Msg: "sample"}, nil
}

func (s *stubQueue) Push(b *Bottle) error {
	return nil
}

func Test_Binn_GetBottle(t *testing.T) {
	nowTime := time.Now()
	now = func() time.Time {
		return nowTime
	}
	defer func() { now = time.Now }()

	q := &stubQueue{}
	itv := time.Duration(1)
	bn := &Binn{
		bq:   q,
		subs: []*Subscription{{id: "example_id", nextTime: now().Add(itv)}},
		itv:  itv,
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
