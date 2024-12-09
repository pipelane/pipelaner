package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pipelane/pipelaner/gen/settings/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server *http.Server
}

func NewMetricsServer(cfg *metrics.Config) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}
	http.Handle(cfg.Path, promhttp.Handler())
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &Server{
		server: server,
	}, nil
}

func (m *Server) Serve(_ context.Context) error {
	return m.server.ListenAndServe()
}

func (m *Server) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
