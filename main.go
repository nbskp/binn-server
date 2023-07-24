package main

import (
	"log"
	"os"
	"time"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/config"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server"
	"golang.org/x/exp/slog"
)

var programLevel = new(slog.LevelVar)

func main() {
	c := config.NewFromEnv()
	q := binn.NewBottleQueue(100, 15*time.Minute)
	bn := binn.New(q, c.SendInterval)
	l := slog.New(logutil.NewCtxHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))
	if c.RunEmitLoop {
		if err := bn.RunEmitLoop(); err != nil {
			l.Error("failed to run emit loop: %v", err)
			return
		}
	}
	srv := server.New(bn, ":8080", l)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
