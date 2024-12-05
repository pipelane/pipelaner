package health

import (
	"context"
	"errors"

	config "github.com/pipelane/pipelaner/gen/pipelaner"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/sources/shared/grpc_server"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	serv *grpc_server.Server
	l    *zerolog.Logger
}

func NewHealthCheck(cfg *config.Pipelaner) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	l, err := logger.NewLoggerWithCfg(cfg.Settings.Logger)
	if err != nil {
		return nil, err
	}
	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host: cfg.Settings.HealthCheck.Host,
		Port: cfg.Settings.HealthCheck.Port,
	}, l)
	return &Server{
		serv: serv,
		l:    l,
	}, nil
}

var (
	errHealthcheckNotInitialized = errors.New("healthcheck server not initialized")
)

func (p *Server) Serve(ctx context.Context) error {
	if p.serv == nil {
		return errHealthcheckNotInitialized
	}
	go func() {
		<-ctx.Done()
		if cErr := p.serv.Stop(); cErr != nil {
			p.l.Error().Err(cErr).Msgf("health check server stop error:")
		}
	}()

	return p.serv.Serve(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	})
}
