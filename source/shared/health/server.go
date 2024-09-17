package health

import (
	"fmt"
	"net"

	"github.com/pipelane/pipelaner"
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
	lis := p.createListener(conf.Host, port)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	go func() {
		err = grpcServer.Serve(lis)
		grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
		if err != nil {
			p.logger.Fatal().Err(err).Msg("Failed run server")
		}
	}()
	return nil
}

func (p *HealthCheck) createListener(host string, port int) net.Listener {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		p.logger.Fatal().Err(err).Msgf("Failed to listen on TCP %s:%d", host, port)
	}
	return tcpListener
}
