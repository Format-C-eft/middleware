package status

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/logger"
)

// NewStatusServer - create status server
func NewStatusServer(cfg *config.Config, isReady *atomic.Value) *http.Server {
	statusAddr := fmt.Sprintf("%s:%v", cfg.Services.Status.Host, cfg.Services.Status.Port)

	mux := http.DefaultServeMux

	mux.HandleFunc(cfg.Services.Status.LivenessPath, livenessHandler)
	mux.HandleFunc(cfg.Services.Status.ReadinessPath, readinessHandler(isReady))
	mux.HandleFunc(cfg.Services.Status.VersionPath, versionHandler(cfg))

	statusServer := &http.Server{
		Addr:    statusAddr,
		Handler: mux,
	}

	return statusServer
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func versionHandler(cfg *config.Config) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		data := map[string]interface{}{
			"name":       cfg.Project.Name,
			"debug":      cfg.Project.Debug,
			"branch":     cfg.Project.Branch,
			"timeBuild":  cfg.Project.TimeBuild,
			"commitHash": cfg.Project.CommitHash,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.ErrorKV(context.Background(), "Service information encoding error", "err", err)
		}
	}
}