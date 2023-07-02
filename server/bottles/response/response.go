package response

import "github.com/nbskp/binn-server/binn"

type Response struct {
	ID        string `json:"id"`
	Msg       string `json:"message"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

func ToResponse(b *binn.Bottle) *Response {
	return &Response{
		ID:        b.ID,
		Msg:       b.Msg,
		Token:     b.Token,
		ExpiredAt: b.ExpiredAt,
	}
}
