package server

import "github.com/nbskp/binn-server/binn"

type bottlesRequest struct {
	ID    string `json:"id"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func (r *bottlesRequest) toBottles() *binn.Bottle {
	return &binn.Bottle{
		ID:    r.ID,
		Msg:   r.Msg,
		Token: r.Token,
	}
}
