package pipelaner

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var AppInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "app_info",
	Help: "Application info",
}, []string{"service_name"})

type MetricsServer struct {
	cfg metricsConfig
}

func NewMetricsServer(cfg Config) (*MetricsServer, error) {
	if !cfg.MetricsEnable {
		return nil, nil
	}
	if cfg.MetricsPort == 0 {
		return nil, fmt.Errorf("metrics port is required")
	}
	prometheus.MustRegister(AppInfo)
	if cfg.MetricsServiceName == "" {
		cfg.MetricsServiceName = "pipelaner"
	}
	AppInfo.WithLabelValues(cfg.MetricsServiceName).Set(1)
	return &MetricsServer{cfg: cfg.metricsConfig}, nil
}

func (m *MetricsServer) Serve() error {
	if !m.cfg.MetricsEnable {
		return nil
	}
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf("%s:%d", m.cfg.MetricsHost, m.cfg.MetricsPort), nil)
}
