package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nbskp/binn-server/binn"
	"github.com/nbskp/binn-server/logutil"
	"github.com/nbskp/binn-server/server"
	"golang.org/x/exp/slog"
)

const (
	envPort         = "PORT"
	envMaxMsgLength = "BINN_MAX_MSG_LENGTH"

	envMySQLDatabase = "MYSQL_Database"
	envMySQLUser     = "MYSQL_USER"
	envMySQLPassword = "MYSQL_PASSWORD"
	envMySQLAddr     = "MYSQL_ADDR"
)

const (
	defaultPort         = "8080"
	defaultMaxMsgLength = 150
)

var programLevel = new(slog.LevelVar)

func main() {
	l := slog.New(logutil.NewCtxHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))

	loc, err := time.LoadLocation("Local")
	if err != nil {
		l.Error(fmt.Sprintf("load timezone: %v", err))
		os.Exit(0)
	}
	dbName := os.Getenv(envMySQLDatabase)
	dbUser := os.Getenv(envMySQLUser)
	dbPassword := os.Getenv(envMySQLPassword)
	dbAddr := os.Getenv(envMySQLAddr)
	c := mysql.Config{
		DBName:    dbName,
		User:      dbUser,
		Passwd:    dbPassword,
		Addr:      dbAddr,
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       loc,
	}
	db, err := sqlx.Connect("mysql", c.FormatDSN())
	if err != nil {
		l.Error(fmt.Sprintf("connect mysql: %v", err))
		os.Exit(0)
	}

	var maxMsgLength int
	if i, err := strconv.Atoi(os.Getenv(envMaxMsgLength)); err != nil {
		l.Warn(fmt.Sprintf("max message length use default value %s", defaultMaxMsgLength))
		maxMsgLength = defaultMaxMsgLength
	} else {
		maxMsgLength = i
	}

	bn := binn.NewBinn(db, maxMsgLength)

	port := os.Getenv(envPort)
	if port == "" {
		l.Warn(fmt.Sprintf("port use default value %s", defaultPort))
		port = defaultPort
	}

	srv := server.New(bn, fmt.Sprintf(":%s", port), l)
	if err := srv.ListenAndServe(); err != nil {
		l.Error(fmt.Sprintf("running server is failed: %v", err))
		os.Exit(0)
	}
}
