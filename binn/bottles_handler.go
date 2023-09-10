package binn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	stateUnavailable = iota
	stateAvailable   = iota
)

type statefulBottle struct {
	bottle *Bottle
	state  int
}

func (sb *statefulBottle) reset() {
	sb.bottle.ExpiredAt = time.Now()
	sb.state = stateAvailable
}

type bottlesHandler struct {
	sbs        []statefulBottle
	size       int
	cnt        int
	expiration time.Duration
}

func NewBottlesHandler(size int, expiration time.Duration) *bottlesHandler {
	sbs := make([]statefulBottle, size)
	for i := 0; i < size; i++ {
		sbs[i].bottle = &Bottle{ID: strconv.Itoa(i), Msg: "", ExpiredAt: time.Now()}
		sbs[i].state = stateAvailable
	}
	return &bottlesHandler{
		sbs:        sbs,
		size:       size,
		cnt:        0,
		expiration: expiration,
	}
}

func (bq *bottlesHandler) Set(ctx context.Context, b *Bottle) error {
	var sb *statefulBottle
	for i := 0; i < bq.size; i++ {
		if bq.sbs[i].bottle.ID == b.ID {
			sb = &bq.sbs[i]
		}
	}
	if sb == nil {
		return NewBinnError(CodeNotFoundBottle, fmt.Sprintf("not found bottle id is %s", b.ID), nil)
	}

	if sb.state == stateAvailable {
		return NewBinnError(CodeUnavailableBottle, fmt.Sprintf("bottle id is %s has not been popped", b.ID), nil)
	}
	if sb.bottle.IsExpired() {
		sb.reset()
		return NewBinnError(CodeExpiredBottle, fmt.Sprintf("bottle id is %s is expired", b.ID), nil)
	}
	sb.bottle.Msg = b.Msg
	sb.reset()
	return nil
}

func (bq *bottlesHandler) Next(ctx context.Context) (*Bottle, error) {
	var sb *statefulBottle
	for i := 0; i < bq.size; i++ {
		if sb_ := &bq.sbs[bq.cnt%bq.size]; sb_.state == stateAvailable || (sb_.state == stateUnavailable && sb_.bottle.IsExpired()) {
			sb = sb_
			bq.cnt++
			break
		} else {
			bq.cnt++
		}
	}
	if sb == nil {
		return nil, nil
	}
	sb.bottle.ExpiredAt = time.Now().Add(bq.expiration)
	sb.state = stateUnavailable
	return sb.bottle, nil
}

type bottlesMySQLHandler struct {
	db         *sqlx.DB
	expiration time.Duration
}

func NewBottlesMySQLHandler(db *sqlx.DB, expiration time.Duration, initSize int) (*bottlesMySQLHandler, error) {
	ctx := context.Background()
	bh := &bottlesMySQLHandler{
		db:         db,
		expiration: expiration,
	}
	if v, err := bh.isNotInitialized(ctx); err != nil {
		return nil, err
	} else if v {
		bh.init(ctx, initSize)
	}
	return bh, nil
}

type bottlesRecord struct {
	ID        string     `db:"id"`
	Msg       string     `db:"msg"`
	ExpiredAt *time.Time `db:"expired_at"`
	Available bool       `db:"available"`
}

func (bh *bottlesMySQLHandler) isNotInitialized(ctx context.Context) (bool, error) {
	rows, err := bh.db.QueryContext(context.TODO(), "SELECT * FROM bottles")
	if err != nil {
		return false, err
	}
	return !rows.Next(), nil
}

func (bh *bottlesMySQLHandler) init(ctx context.Context, size int) error {
	bottlesRecords := make([]*bottlesRecord, size)
	for i := 0; i < size; i++ {
		bottlesRecords[i] = &bottlesRecord{ID: strconv.Itoa(i), Msg: "", ExpiredAt: nil, Available: true}
	}
	_, err := bh.db.NamedExecContext(ctx, "INSERT INTO bottles(id, msg, expired_at, available) VALUES (:id, :msg, :expired_at, :available)", bottlesRecords)
	if err != nil {
		return err
	}
	return nil
}

func (bh *bottlesMySQLHandler) Set(ctx context.Context, b *Bottle) error {
	r, err := bh.db.NamedExecContext(ctx, "UPDATE bottles SET msg=:msg WHERE id=:id AND expired_at>CURRENT_TIMESTAMP AND available=false",
		bottlesRecord{ID: b.ID, Msg: b.Msg, Available: true})
	if err != nil {
		return err
	}
	if n, err := r.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return NewBinnError(CodeNotFoundBottle, fmt.Sprintf("not found bottle id is %s", b.ID), nil)
	}
	return nil
}

func (bh *bottlesMySQLHandler) Next(ctx context.Context) (*Bottle, error) {
	row := bh.db.QueryRowxContext(ctx, "SELECT * FROM bottles WHERE available=true OR (available=false AND expired_at<=?) LIMIT 1", now().Format(time.RFC3339))
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to get a next bottle: %w", err)
	}
	var rc bottlesRecord
	if err := row.StructScan(&rc); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan a next bottle: %w", err)
	}
	expiredAt := time.Now().Add(bh.expiration)
	_, err := bh.db.ExecContext(context.TODO(), "UPDATE bottles SET expired_at=?, available=false WHERE id=?", expiredAt.Format(time.RFC3339), rc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update expired_at and available: %w", err)
	}
	return &Bottle{ID: rc.ID, Msg: rc.Msg, ExpiredAt: time.Now()}, nil
}

type bottlesRedisHandler struct {
	cli  *redis.Client
	exp  time.Duration
	size int
}

func bottleID(n int) string {
	return strconv.Itoa(n)
}

func bottleKey(id string) string {
	return fmt.Sprintf("bottle:%s", id)
}

func bottleShadowKey(id string) string {
	return fmt.Sprintf("bottle:%s.shadow", id)
}

func (h *bottlesRedisHandler) Set(ctx context.Context, b *Bottle) error {
	_, err := h.cli.HGetAll(ctx, bottleShadowKey(b.ID)).Result()
	if err != nil {
		if err == redis.Nil {
			return NewBinnError(CodeNotFoundBottle, "not found the bottle", err)
		}
		return err
	}
	if _, err := h.cli.Del(ctx, bottleShadowKey(b.ID)).Result(); err != nil {
		return err
	}
	if _, err := h.cli.HSet(ctx, bottleKey(b.ID), "msg", b.Msg).Result(); err != nil {
		return err
	}
	return nil
}

func (h *bottlesRedisHandler) Next(ctx context.Context) (*Bottle, error) {
	var b *Bottle
	for i := 0; i < h.size; i++ {
		id := bottleID(i)
		ex, err := h.cli.Exists(ctx, bottleShadowKey(id)).Result()
		if err != nil {
			return nil, err
		}
		if ex == 0 {
			bv, err := h.cli.HGetAll(ctx, bottleKey(id)).Result()
			if err != nil {
				return nil, err
			}
			b = &Bottle{
				ID:  id,
				Msg: bv["msg"],
			}
			break
		}
	}
	if b == nil {
		return nil, nil
	}
	sKey := bottleShadowKey(b.ID)
	if _, err := h.cli.HSet(ctx, sKey, 0, 0).Result(); err != nil {
		return nil, err
	}
	expiredAt := now().Add(h.exp)
	if _, err := h.cli.ExpireAt(ctx, sKey, expiredAt).Result(); err != nil {
		return nil, err
	}
	b.ExpiredAt = expiredAt
	return b, nil
}

func NewBottlesRedisHandler(ctx context.Context, cli *redis.Client, size int, exp time.Duration) (*bottlesRedisHandler, error) {
	ks, err := cli.Keys(ctx, "bottle:*").Result()
	if err != nil {
		return nil, err
	}
	if len(ks) == 0 {
		for i := 0; i < size; i++ {
			_, err := cli.HSet(ctx, bottleKey(bottleID(i)), "msg", "").Result()
			if err != nil {
				return nil, err
			}
		}
	}
	return &bottlesRedisHandler{
		cli:  cli,
		exp:  exp,
		size: size,
	}, nil
}
