package binn

import (
	"context"
	"errors"
)

type subscriptionsHandler struct {
	subscriptions map[string]*Subscription
}

func NewSubscriptionsHandler() *subscriptionsHandler {
	return &subscriptionsHandler{
		subscriptions: map[string]*Subscription{},
	}
}

func (h *subscriptionsHandler) Get(ctx context.Context, id string) (*Subscription, error) {
	sub, ok := h.subscriptions[id]
	if !ok {
		return nil, NewBinnError(CodeNotFoundSubscribedBottle, "not found subscribed a bottle", nil)
	}
	if sub.IsExpired() {
		delete(h.subscriptions, id)
		return nil, NewBinnError(CodeExpiredSubscription, "subscriptions is expired", nil)
	}
	return sub, nil
}

func (h *subscriptionsHandler) Update(ctx context.Context, sub *Subscription) error {
	sub, ok := h.subscriptions[sub.id]
	if !ok {
		return NewBinnError(CodeNotFoundSubscribedBottle, "not found subscribed a bottle", nil)
	}
	if sub.IsExpired() {
		delete(h.subscriptions, sub.id)
		return NewBinnError(CodeExpiredSubscription, "subscriptions is expired", nil)
	}
	h.subscriptions[sub.id] = sub
	return nil
}

func (h *subscriptionsHandler) Add(ctx context.Context, sub *Subscription) error {
	if _, ok := h.subscriptions[sub.id]; ok {
		return errors.New("this subscription is already exists")
	}
	h.subscriptions[sub.id] = sub
	return nil
}
