package bottles

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/logutil"
	"golang.org/x/exp/slog"
	"nhooyr.io/websocket"
)

func WebsocketHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	hf := websocketHandlerFunc(bn, logger)
	return func(w http.ResponseWriter, r *http.Request) {
		hf(w, r)
	}
}

func websocketHandlerFunc(bn *binn.Binn, logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
		defer cancel()

		ctx = c.CloseRead(ctx)

		ch := make(chan struct{}, 0)

		bn.Subscribe(func(b *binn.Bottle) bool {
			select {
			case <-ctx.Done():
				close(ch)
				return false
			default:
			}

			wr, err := c.Writer(ctx, websocket.MessageText)
			if err != nil {
				log.Println(err)
				close(ch)
				return false
			}
			err = json.NewEncoder(wr).Encode(b)
			if err != nil {
				log.Println(err)
				close(ch)
				return false
			}
			logger.Info("emit bottle", logutil.AttrBottle(b))
			wr.Close()
			return true
		})

		select {
		case <-ctx.Done():
		case <-ch:
		}
	})

}
