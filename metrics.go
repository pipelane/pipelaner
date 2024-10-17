package pipelaner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version string // Version Provisioned by ldflags
	Commit  string // Commit Provisioned by ldflags
)

var AppInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "app_info",
	Help: "Application info",
}, []string{"version", "commit", "service_name"})

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
	prometheus.MustRegister(AppInfo)
	if cfg.MetricsServiceName == "" {
		cfg.MetricsServiceName = "pipelaner"
	}
	AppInfo.WithLabelValues(Version, Commit, cfg.MetricsServiceName).Set(1)
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
