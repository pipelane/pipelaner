package pipelaner

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

type MetricsServer struct {
	server *http.Server
}

func NewMetricsServer(cfg *metrics.MetricsConfig) (*MetricsServer, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}
	if cfg.Port == 0 {
		return nil, errors.New("port is required")
	}
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &MetricsServer{
		server: server,
	}, nil
}

func (m *MetricsServer) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		if cErr := m.server.Shutdown(context.Background()); cErr != nil {
			log.Printf("failed to shutdown metrics server: %v", cErr)
		}
	}()
	return m.server.ListenAndServe()
}

func (m *MetricsServer) Close() error {
	return m.server.Close()
}
