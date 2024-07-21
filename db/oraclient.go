package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func NewOracle(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}

	go runScheduler(conn)

	return conn, nil
}

func runScheduler(conn *sql.DB) {
	timer := time.NewTicker(time.Minute * 15)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			ping(conn)
		}
	}
}

func ping(conn *sql.DB) {
	_, err := conn.ExecContext(context.Background(), "SELECT * FROM DUAL")
	if err != nil {
		log.Println("Ping to DB has got an error ::: " + err.Error())
		return
	}
}
