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
	subs map[string]*Subscription

	itv    time.Duration
	subExp time.Duration
}

func NewBinn(itv time.Duration, bq BottleQueue, subExp time.Duration) *Binn {
	return &Binn{
		bq:     bq,
		subs:   map[string]*Subscription{},
		itv:    itv,
		subExp: subExp,
	}
}

type Subscription struct {
	id        string
	expiredAt time.Time

	nextTime  time.Time
	bottleIDs map[string]struct{}
}

func (s *Subscription) IsExpired() bool {
	return now().After(s.expiredAt)
}

func (bn *Binn) Subscribe(id string) {
	bn.subs[id] = &Subscription{
		id:        id,
		expiredAt: now().Add(bn.subExp),
		nextTime:  now().Add(bn.itv),
		bottleIDs: map[string]struct{}{},
	}
}

func (bn *Binn) GetBottle(subID string) (*Bottle, error) {
	if sub, ok := bn.subs[subID]; ok {
		if sub.IsExpired() {
			delete(bn.subs, subID)
			return nil, errors.New("subscription is expired")
		}
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
	return nil, errors.New("not found subscription")
}

func (bn *Binn) Publish(subID string, b *Bottle) error {
	for _, sub := range bn.subs {
		if sub.id == subID {
			if sub.IsExpired() {
				return errors.New("subscription is expired")
			}
			if _, ok := sub.bottleIDs[b.ID]; ok {
				delete(sub.bottleIDs, b.ID)
				return bn.bq.Push(b)
			}
			return errors.New("not found subscribed a bottle")
		}
	}
	return errors.New("not found subscription")
}
