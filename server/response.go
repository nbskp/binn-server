package server

import (
	"time"

	"github.com/nbskp/binn-server/binn"
)

type bottlesResponse struct {
	ID        string    `json:"id"`
	Msg       string    `json:"msg"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func toBottlesResponse(b *binn.Bottle) *bottlesResponse {
	return &bottlesResponse{
		ID:        b.ID,
		Msg:       b.Msg,
		Token:     b.Token,
		ExpiredAt: b.ExpiredAt,
	}
}
