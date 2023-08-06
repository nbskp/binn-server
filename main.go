package main

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nbskp/binn-server/auth"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/config"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server"
	"golang.org/x/exp/slog"
)

var programLevel = new(slog.LevelVar)

func main() {
	c := config.NewFromEnv()

	db, err := sqlx.Connect("mysql", "binn:binn@tcp(mysql:3306)/binn?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	bh, err := binn.NewBottlesMySQLHandler(db, 15*time.Minute, 10)
	if err != nil {
		log.Fatal(err)
	}
	sh := binn.NewSubscriptionsMySQLHandler(db)
	bn := binn.NewBinn(c.SendInterval, bh, sh, c.SubscriptionExpiration)
	l := slog.New(logutil.NewCtxHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))
	auth := auth.NewTokenProvider(10)
	srv := server.New(bn, auth, ":8080", l)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
