package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewPostgresConn(ctx context.Context, dsn string) (*pgx.Conn) {
    conn, err := pgx.Connect(ctx, dsn)
    if err != nil {
        panic("connect with db lost")
    }

    return conn
}