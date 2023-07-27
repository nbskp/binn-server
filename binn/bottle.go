package binn

type Bottle struct {
	ID        string
	Msg       string
	ExpiredAt int64
}

func (b *Bottle) IsExpired() bool {
	return b.ExpiredAt < now().Unix()
}
