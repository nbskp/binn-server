package bottles

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nbskp/binn-server/binn"
)

type request struct {
	ID    string `json:"id"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func (r *request) toBottle() *binn.Bottle {
	return &binn.Bottle{
		ID:    r.ID,
		Msg:   r.Msg,
		Token: r.Token,
	}
}

func HandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	postHf := postHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHf(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
		if err := bn.Publish(reqBody.toBottle()); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
