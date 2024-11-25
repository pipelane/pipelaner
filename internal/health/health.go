package health

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/pipelane/pipelaner/gen/settings/healthcheck"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/source/shared/grpc_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	serv *grpc_server.Server
}

func NewHealthCheck(cfg *healthcheck.HealthcheckConfig) (*Server, error) {
	logger := logger.NewLogger()

	if cfg == nil {
		return nil, errors.New("config is required")
	}

	if cfg.Port == 0 {
		return nil, fmt.Errorf("health check port is required")
	}

	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
	}, &logger)

	return &Server{
		serv: serv,
	}, nil
}

var (
	healthcheckNotInitialized = errors.New("healthcheck server not initialized")
)

func (p *Server) Serve(ctx context.Context) error {
	if p.serv == nil {
		return healthcheckNotInitialized
	}
	go func() {
		<-ctx.Done()
		if cErr := p.serv.Stop(); cErr != nil {
			log.Printf("health check server stop error: %v", cErr)
		}
	}()

	return p.serv.Serve(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	})
}
