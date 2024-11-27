package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"

	ch "github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/pipelane/pipelaner/gen/source/sink"
)

var r = regexp.MustCompile(`(?P<Size>[0-9]+)(?P<Parameter>[a-zA-Z]+)?`)

func ParseSize(text string) (int, error) {
	if len(text) == 0 || text[0] < '0' || text[0] > '9' {
		return 0, errors.New("incorrect size")
	}
	matches := r.FindStringSubmatch(text)

	size, err := strconv.Atoi(matches[r.SubexpIndex("Size")])
	if err != nil {
		return 0, err
	}
	parameter := matches[r.SubexpIndex("Parameter")]

	switch parameter {
	case "GB", "Gb", "gb":
		return size << 30, nil
	case "MB", "Mb", "mb":
		return size << 20, nil
	case "KB", "Kb", "kb":
		return size << 10, nil
	case "B", "b", "":
		return size, nil
	default:
		return 0, fmt.Errorf("unsupported size postfix: %s", parameter)
	}
}

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

type ClientClickhouse struct {
	conn driver.Conn
}

func NewClickhouseClient(ctx context.Context, cfg sink.Clickhouse) (*ClientClickhouse, error) {
	maxCompressionBuffer, err := ParseSize(cfg.GetMaxCompressionBuffer().String())
	if err != nil {
		return nil, err
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.GetAddress()},
		Auth: clickhouse.Auth{
			Database: cfg.GetDatabase(),
			Username: cfg.GetUser(),
			Password: cfg.GetPassword(),
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: cfg.GetEnableDebug(),
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": cfg.GetMaxExecutionTime().GoDuration().Seconds(),
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          cfg.GetDialTimeout().GoDuration(),
		MaxOpenConns:         cfg.GetMaxOpenConns(),
		MaxIdleConns:         cfg.GetMaxIdleConns(),
		ConnMaxLifetime:      cfg.GetCannMaxLifeTime().GoDuration(),
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      cfg.GetBlockBufferSize(),
		MaxCompressionBuffer: maxCompressionBuffer,
	})
	if err != nil {
		return nil, err
	}
	if errs := conn.Ping(ctx); errs != nil {
		return nil, errs
	}
	return &ClientClickhouse{
		conn: conn,
	}, nil
}

func (c *ClientClickhouse) Conn() driver.Conn {
	return c.conn
}
