package binn

import "time"

type Bottle struct {
	ID        string
	Msg       string
	ExpiredAt time.Time
}

func (b *Bottle) IsExpired() bool {
	return b.ExpiredAt.Before(now())
}
