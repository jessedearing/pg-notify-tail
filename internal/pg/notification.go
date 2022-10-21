package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// NotifyOnChannel sends the notification received on a LISTEN postgres topic
// to a Go channel
func NotifyOnChannel(ctx context.Context, conn *pgx.Conn, ch chan<- pgconn.Notification, errch chan<- error) {
	for {
		noti, err := conn.WaitForNotification(ctx)
		if err != nil {
			errch <- err
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		var notiCpy = *noti
		ch <- notiCpy
	}
}
