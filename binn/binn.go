package binn

import (
	"time"
)

type BottleQueue interface {
	Push(*Bottle) error
	Pop() (*Bottle, error)
}

type Binn struct {
	handlers     []handler
	bottleQueue  BottleQueue
	emitInterval time.Duration
}

func New(bottleQueue BottleQueue, emitInterval time.Duration) *Binn {
	bn := &Binn{
		handlers:     []handler{},
		bottleQueue:  bottleQueue,
		emitInterval: emitInterval,
	}
	go bn.emitLoop()
	return bn
}

func (bn *Binn) Run() {
	go bn.emitLoop()
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

func (bn *Binn) emitLoop() {
	for {
		select {
		case <-time.After(bn.emitInterval):
			bn.emit()
		}
	}
}

func (bn *Binn) emit() {
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
