package metrics

import (
	"fmt"
	"net/http"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMetricsServer - new metrics server
func NewMetricsServer(cfg *config.Config) *http.Server {

	addr := fmt.Sprintf("%s:%d", cfg.Services.Metrics.Host, cfg.Services.Metrics.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Services.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return metricsServer
}
