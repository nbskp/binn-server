package request

import "github.com/nbskp/binn-server/binn"

type Request struct {
	ID    string `json:"id"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func (r *Request) ToBottle() *binn.Bottle {
	return &binn.Bottle{
		ID:    r.ID,
		Msg:   r.Msg,
		Token: r.Token,
	}
}
