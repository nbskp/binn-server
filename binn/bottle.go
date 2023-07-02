package binn

import (
	"time"
)

type Bottle struct {
	ID        string
	Msg       string
	Token     string
	ExpiredAt int64
}

var now = time.Now

func (b *Bottle) IsExpired() bool {
	return b.ExpiredAt < now().Unix()
}
