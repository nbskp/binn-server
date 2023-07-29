package server

import "github.com/nbskp/binn-server/binn"

type bottlesRequest struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

func (r *bottlesRequest) toBottles() *binn.Bottle {
	return &binn.Bottle{
		ID:  r.ID,
		Msg: r.Msg,
	}
}
