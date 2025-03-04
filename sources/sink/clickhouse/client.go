/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package clickhouse

import (
	"context"
	"time"

	ch "github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/pipelane/pipelaner/gen/source/sink"
)

type Client struct {
	conn *chpool.Pool
}

func NewClickhouseClient(ctx context.Context, cfg sink.Clickhouse) (*Client, error) {
	conn, err := chpool.Dial(ctx, chpool.Options{
		ClientOptions: ch.Options{
			Address:          cfg.GetCredentials().Address,
			Database:         cfg.GetCredentials().Database,
			User:             cfg.GetCredentials().User,
			Password:         cfg.GetCredentials().Password,
			Compression:      ch.CompressionLZ4,
			DialTimeout:      5 * time.Second,
			HandshakeTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) GetConn(ctx context.Context) (*chpool.Client, error) {
	conn, err := c.conn.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
