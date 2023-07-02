package main

import (
	"log"
	"time"

	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/config"
	"github.com/nbskp/binn-server/server"
)

func main() {
	c := config.NewFromEnv()
	q := binn.NewBottleQueue(100, 15*time.Minute)
	bn := binn.New(q, c.SendInterval)
	l := log.Default()
	srv := server.New(bn, ":8080", l)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
