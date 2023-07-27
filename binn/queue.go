package binn

import (
	"errors"
	"fmt"
	"strconv"
	"time"
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
	sb.bottle.ExpiredAt = 0
	sb.state = stateAvailable
}

type bottleQueue struct {
	sbs        []statefulBottle
	size       int
	cnt        int
	expiration time.Duration
}

func NewBottleQueue(size int, expiration time.Duration) *bottleQueue {
	sbs := make([]statefulBottle, size)
	for i := 0; i < size; i++ {
		sbs[i].bottle = &Bottle{ID: strconv.Itoa(i), Msg: "", ExpiredAt: 0}
		sbs[i].state = stateAvailable
	}
	return &bottleQueue{
		sbs:        sbs,
		size:       size,
		cnt:        0,
		expiration: expiration,
	}
}

func (bq *bottleQueue) Push(b *Bottle) error {
	var sb *statefulBottle
	for i := 0; i < bq.size; i++ {
		if bq.sbs[i].bottle.ID == b.ID {
			sb = &bq.sbs[i]
		}
	}
	if sb == nil {
		return errors.New(fmt.Sprintf("not found bottle id is %s", b.ID))
	}

	if sb.state == stateAvailable {
		return errors.New(fmt.Sprintf("bottle id is %s has not been popped", b.ID))
	}
	if sb.bottle.IsExpired() {
		sb.reset()
		return errors.New(fmt.Sprintf("bottle id is %s is expired", b.ID))
	}
	sb.bottle.Msg = b.Msg
	sb.reset()
	return nil
}

func (bq *bottleQueue) Pop() (*Bottle, error) {
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
	sb.bottle.ExpiredAt = time.Now().Add(bq.expiration).Unix()
	sb.state = stateUnavailable
	return sb.bottle, nil
}
