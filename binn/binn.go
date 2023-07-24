package binn

import (
	"context"
	"errors"
	"time"
)

type BottleQueue interface {
	Push(*Bottle) error
	Pop() (*Bottle, error)
}

type Binn struct {
	handlers       []handler
	bottleQueue    BottleQueue
	emitInterval   time.Duration
	cancelEmitLoop func()
}

func New(bottleQueue BottleQueue, emitInterval time.Duration) *Binn {
	bn := &Binn{
		handlers:     []handler{},
		bottleQueue:  bottleQueue,
		emitInterval: emitInterval,
	}
	return bn
}

func (bn *Binn) Publish(b *Bottle) error {
	if err := bn.bottleQueue.Push(b); err != nil {
		return err
	}
	return nil
}

type handler func(*Bottle) bool

func (bn *Binn) Subscribe(h handler) error {
	bn.handlers = append(bn.handlers, h)
	return nil
}

func (bn *Binn) RunEmitLoop() error {
	if bn.cancelEmitLoop != nil {
		return errors.New("emit-loop is already running")
	}
	var ctx context.Context
	ctx, bn.cancelEmitLoop = context.WithCancel(context.Background())
	go bn.emitLoop(ctx)
	return nil
}

func (bn *Binn) StopEmitLoop() error {
	if bn.cancelEmitLoop == nil {
		return errors.New("emit-loop isn't running")
	}
	bn.cancelEmitLoop()
	bn.cancelEmitLoop = nil
	return nil
}

func (bn *Binn) emitLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(bn.emitInterval):
			bn.Emit()
		}
	}
}

func (bn *Binn) Emit() {
	if lh := len(bn.handlers); lh == 0 {
		return
	}
	b, err := bn.bottleQueue.Pop()
	if err != nil {
		return
	}
	fn := bn.handlers[0]
	bn.handlers = bn.handlers[1:]

	if ok := fn(b); !ok {
		_ = bn.bottleQueue.Push(b)
		return
	}
	bn.handlers = append(bn.handlers, fn)
}
