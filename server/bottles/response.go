package bottles

import "github.com/nbskp/binn-server/binn"

type Response struct {
	ID        string `json:"id"`
	Msg       string `json:"message"`
	ExpiredAt int64  `json:"expired_at"`
}

func ToResponse(b *binn.Bottle) *Response {
	return &Response{
		ID:        b.ID,
		Msg:       b.Msg,
		ExpiredAt: b.ExpiredAt,
	}
}
