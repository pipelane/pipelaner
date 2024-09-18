package health

import (
	"fmt"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/source/shared/grpc_server"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Config struct {
	Host              string `pipelaner:"host"`
	Port              *int   `pipelane:"port"`
	EnableHealthCheck bool   `pipelane:"enable_health_check"`
}

//nolint:revive
type HealthCheck struct {
	logger zerolog.Logger
	cfg    *pipelaner.BaseLaneConfig
}

func NewHealthCheck(logger zerolog.Logger) *HealthCheck {
	return &HealthCheck{
		logger: logger,
	}
}

func (p *HealthCheck) Init(ctx *pipelaner.Context) error {
	p.cfg = ctx.LaneItem().Config()
	p.logger = pipelaner.NewLogger()
	conf := &Config{}
	err := p.cfg.ParseExtended(conf)
	if err != nil {
		return err
	}

	if !conf.EnableHealthCheck {
		return nil
	}

	if conf.Host == "" {
		return fmt.Errorf("host is required")
	}

	port := 84
	if conf.Port != nil {
		port = *conf.Port
	}

	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host: conf.Host,
		Port: port,
	}, p.logger)

	serv.Serve(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	})

	return nil
}
