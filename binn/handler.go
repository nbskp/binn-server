package binn

import (
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

type bottlesHandler struct {
	sbs        []statefulBottle
	size       int
	cnt        int
	expiration time.Duration
}

func NewBottlesHandler(size int, expiration time.Duration) *bottlesHandler {
	sbs := make([]statefulBottle, size)
	for i := 0; i < size; i++ {
		sbs[i].bottle = &Bottle{ID: strconv.Itoa(i), Msg: "", ExpiredAt: 0}
		sbs[i].state = stateAvailable
	}
	return &bottlesHandler{
		sbs:        sbs,
		size:       size,
		cnt:        0,
		expiration: expiration,
	}
}

func (bq *bottlesHandler) Set(b *Bottle) error {
	var sb *statefulBottle
	for i := 0; i < bq.size; i++ {
		if bq.sbs[i].bottle.ID == b.ID {
			sb = &bq.sbs[i]
		}
	}
	if sb == nil {
		return NewBinnError(CodeNotFoundBottle, fmt.Sprintf("not found bottle id is %s", b.ID), nil)
	}

	if sb.state == stateAvailable {
		return NewBinnError(CodeUnavailableBottle, fmt.Sprintf("bottle id is %s has not been popped", b.ID), nil)
	}
	if sb.bottle.IsExpired() {
		sb.reset()
		return NewBinnError(CodeExpiredBottle, fmt.Sprintf("bottle id is %s is expired", b.ID), nil)
	}
	sb.bottle.Msg = b.Msg
	sb.reset()
	return nil
}

func (bq *bottlesHandler) Next() (*Bottle, error) {
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
