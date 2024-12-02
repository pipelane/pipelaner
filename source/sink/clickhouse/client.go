package clickhouse

import (
	"context"
	"time"

	ch "github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/pipelane/pipelaner/gen/source/sink"
)

type LowLevelClickhouseClient struct {
	conn *chpool.Pool
}

func NewLowLevelClickhouseClient(ctx context.Context, cfg sink.Clickhouse) (*LowLevelClickhouseClient, error) {
	conn, err := chpool.Dial(ctx, chpool.Options{

		ClientOptions: ch.Options{
			Address:          cfg.GetAddress(),
			Database:         cfg.GetDatabase(),
			User:             cfg.GetUser(),
			Password:         cfg.GetPassword(),
			Compression:      ch.CompressionLZ4,
			DialTimeout:      5 * time.Second,
			HandshakeTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		return nil, err
	}

	return &LowLevelClickhouseClient{
		conn: conn,
	}, nil
}

func (c *LowLevelClickhouseClient) GetConn(ctx context.Context) (*chpool.Client, error) {
	conn, err := c.conn.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
