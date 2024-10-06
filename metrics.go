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

func StartMetricsServer(cfg metricsConfig) error {
	if !cfg.MetricsEnable {
		return nil
	}
	if cfg.MetricsPort == 0 {
		return fmt.Errorf("metrics port is required")
	}
	prometheus.MustRegister(AppInfo)
	if cfg.MetricsServiceName == "" {
		cfg.MetricsServiceName = "pipelaner"
	}
	AppInfo.WithLabelValues(cfg.MetricsServiceName).Set(1)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.MetricsHost, cfg.MetricsPort), nil)
}
