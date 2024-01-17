package binn

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Binn struct {
	db           *sqlx.DB
	maxMsgLength int
}

func NewBinn(db *sqlx.DB, maxMsgLength int) *Binn {
	return &Binn{
		db:           db,
		maxMsgLength: maxMsgLength,
	}
}

func (bn *Binn) Get(ctx context.Context) (*Bottle, error) {
	return nil, nil
}

func (bn *Binn) Set(ctx context.Context, b *Bottle) error {
	return nil
}
