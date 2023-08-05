package binn

import (
	"context"
	"time"
)

type BottlesHandler interface {
	Set(context.Context, *Bottle) error
	Next(context.Context) (*Bottle, error)
}

type Binn struct {
	bh   BottlesHandler
	subs map[string]*Subscription

	itv    time.Duration
	subExp time.Duration
}

func NewBinn(itv time.Duration, bq BottlesHandler, subExp time.Duration) *Binn {
	return &Binn{
		bh:     bq,
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

func (bn *Binn) Subscribe(ctx context.Context, id string) {
	bn.subs[id] = &Subscription{
		id:        id,
		expiredAt: now().Add(bn.subExp),
		nextTime:  now().Add(bn.itv),
		bottleIDs: map[string]struct{}{},
	}
}

func (bn *Binn) GetBottle(ctx context.Context, subID string) (*Bottle, error) {
	if sub, ok := bn.subs[subID]; ok {
		if sub.IsExpired() {
			delete(bn.subs, subID)
			return nil, NewBinnError(CodeExpiredSubscription, "subscriptions is expired", nil)
		}
		if sub.nextTime.After(now()) {
			return nil, nil
		}
		b, err := bn.bh.Next(ctx)
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
	return nil, NewBinnError(CodeNotFoundSubscription, "not found subscription", nil)
}

func (bn *Binn) Publish(ctx context.Context, subID string, b *Bottle) error {
	for _, sub := range bn.subs {
		if sub.id == subID {
			if sub.IsExpired() {
				return NewBinnError(CodeExpiredSubscription, "subscriptions is expired", nil)
			}
			if _, ok := sub.bottleIDs[b.ID]; ok {
				delete(sub.bottleIDs, b.ID)
				return bn.bh.Set(ctx, b)
			}
			return NewBinnError(CodeNotFoundSubscribedBottle, "not found subscribed a bottle", nil)
		}
	}
	return NewBinnError(CodeNotFoundSubscription, "not found subscription", nil)
}
