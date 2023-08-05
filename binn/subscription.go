package binn

import "time"

type Subscription struct {
	id        string
	expiredAt time.Time

	nextTime  time.Time
	bottleIDs []string
}

func (s *Subscription) IsExpired() bool {
	return now().After(s.expiredAt)
}
