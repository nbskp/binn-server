package server

import (
	"time"

	"github.com/nbskp/binn-server/binn"
)

type bottlesResponse struct {
	ID        string    `json:"id"`
	Msg       string    `json:"message"`
	ExpiredAt time.Time `json:"expired_at"`
}

func toBottlesResponse(b *binn.Bottle) *bottlesResponse {
	return &bottlesResponse{
		ID:        b.ID,
		Msg:       b.Msg,
		ExpiredAt: b.ExpiredAt,
	}
}

type subscribeBottlesResponse struct {
	Token string `json:"token"`
}

func newSubscribeBottlesResponse(token string) *subscribeBottlesResponse {
	return &subscribeBottlesResponse{Token: token}
}
