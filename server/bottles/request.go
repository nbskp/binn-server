package bottles

import "github.com/nbskp/binn-server/binn"

type Request struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

func (r *Request) ToBottle() *binn.Bottle {
	return &binn.Bottle{
		ID:  r.ID,
		Msg: r.Msg,
	}
}
