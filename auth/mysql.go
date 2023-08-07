package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type tokenMySQLProvider struct {
	db   *sqlx.DB
	tLen int
}

func NewTokenMySQLProvider(db *sqlx.DB, tLen int) *tokenMySQLProvider {
	return &tokenMySQLProvider{
		db:   db,
		tLen: tLen,
	}
}

func (p *tokenMySQLProvider) Issue(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return "", errors.New("timeout issue token")
		default:
		}
		token, err := generateRandomStr(uint32(p.tLen))
		if err != nil {
			return "", err
		}
		if err := p.db.QueryRowxContext(ctx, "SELECT * FROM tokens WHERE id = ?", token).StructScan(&struct{ ID string }{}); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err := p.db.ExecContext(ctx, "INSERT INTO tokens (id) VALUE (?)", token)
				if err != nil {
					return "", nil
				}
				return token, nil
			} else {
				return "", err
			}
		}
	}
}

func (p *tokenMySQLProvider) Authorize(ctx context.Context, token string) (string, bool, error) {
	if err := p.db.QueryRowxContext(ctx, "SELECT * FROM tokens WHERE id = ?", token).StructScan(&struct{ ID string }{}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		} else {
			return "", false, err
		}
	}
	return token, true, nil
}
