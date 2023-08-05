package binn

import (
	"context"
	"time"
)

type BottlesHandler interface {
	Set(context.Context, *Bottle) error
	Next(context.Context) (*Bottle, error)
}

type SubscriptionsHandler interface {
	Get(context.Context, string) (*Subscription, error)
	Add(context.Context, *Subscription) error
	Update(context.Context, *Subscription) error
}

type Binn struct {
	bh BottlesHandler
	sh SubscriptionsHandler

	itv    time.Duration
	subExp time.Duration
}

func NewBinn(itv time.Duration, bh BottlesHandler, sh SubscriptionsHandler, subExp time.Duration) *Binn {
	return &Binn{
		bh:     bh,
		sh:     sh,
		itv:    itv,
		subExp: subExp,
	}
}

func (bn *Binn) Subscribe(ctx context.Context, subID string) error {
	return bn.sh.Add(ctx, &Subscription{
		id:        subID,
		expiredAt: now().Add(bn.subExp),
		nextTime:  now().Add(bn.itv),
		bottleIDs: []string{}},
	)
}

func (bn *Binn) GetBottle(ctx context.Context, subID string) (*Bottle, error) {
	sub, err := bn.sh.Get(ctx, subID)
	if err != nil {
		return nil, err
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
	bn.sh.Update(ctx, &Subscription{
		id:        sub.id,
		expiredAt: sub.expiredAt,
		nextTime:  time.Now().Add(bn.itv),
		bottleIDs: append(sub.bottleIDs, b.ID),
	})
	return b, nil
}

func (bn *Binn) Publish(ctx context.Context, subID string, b *Bottle) error {
	sub, err := bn.sh.Get(ctx, subID)
	if err != nil {
		return err
	}
	isFound := false
	bottleIDs := []string{}
	for _, bottleID := range sub.bottleIDs {
		if bottleID == b.ID {
			isFound = true
		} else {
			bottleIDs = append(bottleIDs, bottleID)
		}
	}
	if isFound {
		bn.sh.Update(ctx, &Subscription{
			id:        sub.id,
			expiredAt: sub.expiredAt,
			nextTime:  sub.nextTime,
			bottleIDs: bottleIDs,
		})
		return nil
	}
	return NewBinnError(CodeNotFoundSubscription, "not found subscription", nil)
}
