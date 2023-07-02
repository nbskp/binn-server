package bottles

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles/request"
)

func NewBottlesMux(bn *binn.Binn, logger *log.Logger) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/", bottlesHandlerFunc(bn, logger))
	r.HandleFunc("/stream", StreamHandlerFunc(bn, logger))
	r.HandleFunc("/ws", WebsocketHandlerFunc(bn, logger))
	return r
}

func bottlesHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	postHf := postBottlesHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHf(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postBottlesHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody request.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}
		if err := bn.Publish(reqBody.ToBottle()); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
