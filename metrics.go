package pipelaner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	cfg metricsConfig
}

func NewMetricsServer(cfg Config) (*MetricsServer, error) {
	if !cfg.MetricsEnable {
		return nil, nil //nolint:nilnil
	}
	if cfg.MetricsPort == 0 {
		return nil, fmt.Errorf("metrics port is required")
	}
	return &MetricsServer{cfg: cfg.metricsConfig}, nil
}

func (m *MetricsServer) Serve() error {
	if !m.cfg.MetricsEnable {
		return nil
	}
	http.Handle("/metrics", promhttp.Handler())
	s := http.Server{
		Addr:              fmt.Sprintf("%s:%d", m.cfg.MetricsHost, m.cfg.MetricsPort),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s.ListenAndServe()
}
