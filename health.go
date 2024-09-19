package pipelaner

import (
	"fmt"

	"github.com/pipelane/pipelaner/source/shared/grpc_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//nolint:revive
type HealthCheck struct {
	serv *grpc_server.Server
}

func NewHealthCheck(conf healthCheckConfig) (*HealthCheck, error) {
	logger := NewLogger()

	if !conf.EnableHealthCheck {
		return nil, nil
	}

	if conf.Host == "" {
		return nil, fmt.Errorf("host is required")
	}

	port := 84
	if conf.Port != nil {
		port = *conf.Port
	}

	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host: conf.Host,
		Port: port,
	}, logger)

	return &HealthCheck{
		serv: serv,
	}, nil
}

func (p *HealthCheck) Serve() {
	p.serv.Serve(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	})
}
