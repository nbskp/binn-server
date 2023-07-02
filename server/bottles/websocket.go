package bottles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/server/bottles/response"
	"golang.org/x/net/websocket"
)

func WebsocketHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	hf := websocketHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		hf(w, r)
	}
}

func websocketHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(func(ws *websocket.Conn) {
			enc := json.NewEncoder(ws)
			closed := make(chan struct{}, 0)
			bn.Subscribe(func(b *binn.Bottle) bool {
				select {
				case <-r.Context().Done():
					close(closed)
					return false
				default:
				}
				if err := enc.Encode(response.ToResponse(b)); err != nil {
					fmt.Println(err)
					close(closed)
					return false
				}
				return true
			})
			select {
			case <-r.Context().Done():
			case <-closed:
			}
			return
		}).ServeHTTP(w, r)
	}
}
