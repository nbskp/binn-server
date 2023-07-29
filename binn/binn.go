package binn

import (
	"errors"
	"time"
)

type BottleQueue interface {
	Push(*Bottle) error
	Pop() (*Bottle, error)
}

type Binn struct {
	bq   BottleQueue
	subs []*Subscription

	itv time.Duration
}

func NewBinn(itv time.Duration, bq BottleQueue) *Binn {
	return &Binn{
		bq:   bq,
		subs: []*Subscription{},
		itv:  itv,
	}
}

type Subscription struct {
	id       string
	nextTime time.Time

	bottleIDs map[string]struct{}
}

func (bn *Binn) Subscribe(id string) {
	bn.subs = append(bn.subs, &Subscription{id: id, nextTime: now().Add(bn.itv), bottleIDs: map[string]struct{}{}})
}

func (bn *Binn) GetBottle(subID string) (*Bottle, error) {
	for _, sub := range bn.subs {
		if sub.id == subID {
			if sub.nextTime.After(now()) {
				return nil, nil
			}
			b, err := bn.bq.Pop()
			if err != nil {
				return nil, err
			}
			if b == nil {
				return nil, nil
			}
			sub.nextTime = time.Now().Add(bn.itv)
			sub.bottleIDs[b.ID] = struct{}{}
			return b, nil
		}
	}
	return nil, errors.New("not found subscription")
}

func (bn *Binn) Publish(subID string, b *Bottle) error {
	for _, sub := range bn.subs {
		if sub.id == subID {
			if _, ok := sub.bottleIDs[b.ID]; ok {
				delete(sub.bottleIDs, b.ID)
				return bn.bq.Push(b)
			}
			return errors.New("not found subscribed a bottle")
		}
	}
	return errors.New("not found subscription")
}
