package metrics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pipelane/pipelaner/gen/settings/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsNotInitializedErr = errors.New("metrics server not initialized")
)

type Server struct {
	server *http.Server
}

func NewMetricsServer(cfg *metrics.MetricsConfig) (*Server, error) {
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

func (m *Server) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		if cErr := m.server.Shutdown(context.Background()); cErr != nil {
			log.Printf("failed to shutdown metrics server: %v", cErr)
		}
	}()
	return m.server.ListenAndServe()
}

func (m *Server) Close() error {
	return m.server.Close()
}
