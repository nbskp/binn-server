package binn

import (
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

func Test_Binn(t *testing.T) {
	q := &stubQueue{}
	bn := &Binn{
		bottleQueue: q,
	}

	ch := make(chan struct{}, 1)
	assertFunc := func(b *Bottle) {
		assert.Equal(t, "sample", b.Msg)
		close(ch)
	}

	bn.Subscribe(func(b *Bottle) bool {
		assertFunc(b)
		return true
	})

	bn.emit()
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "failed to finish")
	}
}
