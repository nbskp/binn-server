package binn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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

type subscriptionsMySQLHandler struct {
	db *sqlx.DB
}

type subscriptionsRecord struct {
	ID        string    `db:"id"`
	ExpiredAt time.Time `db:"expired_at"`
	NextTime  time.Time `db:"next_time"`
	BottleID  string    `db:"bottle_id"`
}

type subscribedBottlesRecord struct {
	BottleID       string `db:"bottle_id"`
	SubscriptionID string `db:"subscription_id"`
}

func (sh *subscriptionsMySQLHandler) Get(ctx context.Context, id string) (*Subscription, error) {
	rows, err := sh.db.QueryxContext(ctx, "SELECT s.id AS id, s.expired_at AS expired_at, s.next_time AS next_time, sb.bottle_id AS bottle_id FROM subscriptions AS s JOIN subscribed_bottles AS sb ON s.id = sb.subscription_id WHERE s.id=?", id)
	if err != nil {
		return nil, err
	}
	sub := &Subscription{}
	if rows.Next() {
		var r subscriptionsRecord
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		sub.id = r.ID
		sub.expiredAt = r.ExpiredAt
		sub.nextTime = r.NextTime
		sub.bottleIDs = []string{r.BottleID}
	}
	for rows.Next() {
		var r subscriptionsRecord
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		sub.bottleIDs = append(sub.bottleIDs, r.BottleID)
	}
	return sub, nil
}

func (sh *subscriptionsMySQLHandler) Update(ctx context.Context, sub *Subscription) error {
	_, err := sh.db.ExecContext(ctx, "UPDATE subscriptions SET expired_at=?, next_time=? WHERE id=?",
		sub.expiredAt.Format(time.RFC3339), sub.nextTime.Format(time.RFC3339), sub.id)
	fmt.Println(sub.expiredAt.Format(time.RFC3339), sub.nextTime.Format(time.RFC3339))
	if err != nil {
		return err
	}
	rs := make([]*subscribedBottlesRecord, len(sub.bottleIDs))
	for i := 0; i < len(sub.bottleIDs); i++ {
		rs[i] = &subscribedBottlesRecord{BottleID: sub.bottleIDs[i], SubscriptionID: sub.id}
	}
	_, err = sh.db.NamedExecContext(ctx, "INSERT IGNORE subscribed_bottles (bottle_id, subscription_id) VALUES (:bottle_id, :subscription_id)", rs)
	if err != nil {
		return err
	}
	return nil
}

func (sh *subscriptionsMySQLHandler) Add(ctx context.Context, sub *Subscription) error {
	_, err := sh.db.ExecContext(ctx, "INSERT INTO subscriptions (id, expired_at, next_time) VALUES (?, ?, ?)",
		sub.id, sub.expiredAt.Format(time.RFC3339), sub.nextTime.Format(time.RFC3339))
	if err != nil {
		return err
	}

	rs := make([]*subscribedBottlesRecord, len(sub.bottleIDs))
	for i := 0; i < len(sub.bottleIDs); i++ {
		rs[i] = &subscribedBottlesRecord{BottleID: sub.bottleIDs[i], SubscriptionID: sub.id}
	}
	_, err = sh.db.NamedExecContext(ctx, "INSERT IGNORE subscribed_bottles (bottle_id, subscription_id) VALUES (:bottle_id, :subscription_id)", rs)
	if err != nil {
		return err
	}

	return nil
}

func NewSubscriptionsMySQLHandler(db *sqlx.DB) *subscriptionsMySQLHandler {
	return &subscriptionsMySQLHandler{db: db}
}

type subscriptionsRedisHandler struct {
	cli *redis.Client
	exp time.Duration
}

func subscriptionKey(id string) string {
	return fmt.Sprintf("subscription:%s", id)
}

func subscriptionToHashFields(s *Subscription) []interface{} {
	fs := make([]interface{}, 4)
	fs = append(fs, "next_time", s.nextTime)
	fs = append(fs, "bottle_ids", strings.Join(s.bottleIDs, ","))
	return fs
}

func mapToSubscription(m map[string]string) (*Subscription, error) {
	nt, err := time.Parse(time.RFC3339, m["next_time"])
	if err != nil {
		return nil, err
	}
	return &Subscription{
		nextTime:  nt,
		bottleIDs: strings.Split(m["bottle_ids"], ","),
	}, nil
}

func (sh *subscriptionsRedisHandler) Get(ctx context.Context, id string) (*Subscription, error) {
	vs, err := sh.cli.HGetAll(ctx, subscriptionKey(id)).Result()
	if err != nil {
		return nil, err
	}
	// https://github.com/redis/go-redis/issues/1668
	if len(vs) == 0 {
		return nil, nil
	}
	s, err := mapToSubscription(vs)
	if err != nil {
		return nil, err
	}
	s.id = id
	return s, nil
}

func (sh *subscriptionsRedisHandler) Update(ctx context.Context, sub *Subscription) error {
	fs := subscriptionToHashFields(sub)
	_, err := sh.cli.HSet(ctx, subscriptionKey(sub.id), fs...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (sh *subscriptionsRedisHandler) Add(ctx context.Context, sub *Subscription) error {
	fs := subscriptionToHashFields(sub)
	_, err := sh.cli.HSet(ctx, subscriptionKey(sub.id), fs...).Result()
	if err != nil {
		return err
	}
	_, err = sh.cli.ExpireAt(ctx, subscriptionKey(sub.id), now().Add(sh.exp)).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewSubscriptionsRedisHandler(cli *redis.Client, exp time.Duration) *subscriptionsRedisHandler {
	return &subscriptionsRedisHandler{cli: cli, exp: exp}
}
