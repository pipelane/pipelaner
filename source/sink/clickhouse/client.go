package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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

type ClientClickhouse struct {
	conn driver.Conn
}

func NewClickhouseClient(ctx context.Context, cfg Config) (*ClientClickhouse, error) {
	maxCompressionBuffer, err := ParseSize(cfg.MaxCompressionBuffer)
	if err != nil {
		return nil, err
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Address},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: cfg.EnableDebug,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": cfg.MaxExecutionTime.Seconds(),
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          cfg.DialTimeout,
		MaxOpenConns:         cfg.MaxOpenConns,
		MaxIdleConns:         cfg.MaxIdleConns,
		ConnMaxLifetime:      cfg.ConnMaxLifetime,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      cfg.BlockBufferSize,
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
