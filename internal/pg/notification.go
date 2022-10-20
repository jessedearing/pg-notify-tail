package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// NotifyOnChannel sends the notification received on a LISTEN postgres topic
// to a Go channel
func NotifyOnChannel(ctx context.Context, conn *pgx.Conn, ch chan<- pgconn.Notification) {
	for {
		noti, err := conn.WaitForNotification(ctx)
		if err != nil {
			close(ch)
			return
		}

		var notiCpy = *noti
		ch <- notiCpy
	}
}
