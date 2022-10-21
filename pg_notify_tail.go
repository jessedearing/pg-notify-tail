package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jessedearing/pg-notify-tail/internal/config"
	"github.com/jessedearing/pg-notify-tail/internal/pg"
)

var (
	pgurl string
)

func main() {
	cfg := config.Config{}
	flag.StringVar(&cfg.PostgresURL, "pgurl", "", "Postgres URL in the form of 'postgres://user:password@host:port/db'")
	flag.StringVar(&cfg.Channel, "channel", "", "The name of the channel to listen on")
	flag.Parse()
	err := cfg.Validate()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err != nil {
		ExitWithError(err)
		return
	}

	conn, err := pgx.Connect(ctx, cfg.PostgresURL)
	if err != nil {
		ExitWithError(err)
		return
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("listen %s", cfg.Channel))
	if err != nil {
		ExitWithError(err)
		return
	}

	var notiChan = make(chan pgconn.Notification)
	var errChan = make(chan error)
	go pg.NotifyOnChannel(ctx, conn, notiChan, errChan)

	for {
		select {
		case n := <-notiChan:
			fmt.Fprintln(os.Stdout, n.Payload)
		case err = <-errChan:
			cancel()
			ExitWithError(err)
		case <-ctx.Done():
			return
		}
	}
}

func ExitWithError(err error) {
	fmt.Fprintln(os.Stderr, "ERROR", err.Error())
	os.Exit(1)
}
