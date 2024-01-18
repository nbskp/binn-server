package binn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Binn struct {
	db           *sqlx.DB
	tokenLength  int
	maxMsgLength int
}

type bottleRecord struct {
	ID         string     `db:"id"`
	Msg        string     `db:"msg"`
	Token      *string    `db:"token"`
	Expiration string     `db:"expiration"`
	ExpiredAt  *time.Time `db:"expired_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func NewBinn(db *sqlx.DB, tokenLength, maxMsgLength int) *Binn {
	return &Binn{
		db:           db,
		tokenLength:  tokenLength,
		maxMsgLength: maxMsgLength,
	}
}

func (bn *Binn) Get(ctx context.Context) (*Bottle, error) {
	row := bn.db.QueryRowxContext(ctx, "SELECT * FROM bottles WHERE expired_at<=CURRENT_TIMESTAMP OR expired_at IS NULL LIMIT 1")
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("get available a bottle: %v", err)
	}
	var rc bottleRecord
	if err := row.StructScan(&rc); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("no ")
			return nil, nil
		}
	}

	tok, err := generateToken(uint32(bn.tokenLength))
	if err != nil {
		return nil, fmt.Errorf("generate a token: %v", err)
	}
	exp, err := time.ParseDuration(rc.Expiration)
	if err != nil {
		return nil, fmt.Errorf("parse expiration: %v", err)
	}
	expiredAt := time.Now().Add(exp)
	_, err = bn.db.NamedExecContext(ctx, "UPDATE bottles SET token=:token, expired_at=:expired_at WHERE id=:id",
		bottleRecord{ID: rc.ID, Token: &tok, ExpiredAt: &expiredAt})
	if err != nil {
		return nil, fmt.Errorf("update a token and expired_at: %v", err)
	}
	return &Bottle{ID: rc.ID, Msg: rc.Msg, Token: tok, ExpiredAt: expiredAt}, nil
}

func (bn *Binn) Set(ctx context.Context, b *Bottle) error {
	r, err := bn.db.NamedExecContext(ctx, "UPDATE bottles SET msg=:msg, token=NULL WHERE id=:id AND expired_at>CURRENT_TIMESTAMP AND token=:token",
		bottleRecord{ID: b.ID, Msg: b.Msg, Token: &b.Token})
	if err != nil {
		return fmt.Errorf("set a bottle: %v", err)
	}
	if n, err := r.RowsAffected(); err != nil {
		return fmt.Errorf("get row affections: %v", err)
	} else if n == 0 {
		return NewBinnError(CodeNotFoundBottle, fmt.Sprintf("not found valid bottle ", b.ID), nil)
	}
	return nil
}
