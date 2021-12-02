package servers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/logger"
	metric "github.com/Format-C-eft/middleware/internal/servers/metrics"
	"github.com/Format-C-eft/middleware/internal/servers/rest"
	"github.com/Format-C-eft/middleware/internal/servers/status"
)

// Start - start rest server
func Start(ctx context.Context) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfigInstance()

	metricsServer := metric.NewMetricsServer(&cfg)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Metrics server is running on %s:%v", cfg.Services.Metrics.Host, cfg.Services.Metrics.Port))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WarnKV(ctx, "Failed running metrics server", "err", err)
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := status.NewStatusServer(&cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", cfg.Services.Status.Host, cfg.Services.Status.Port)
		logger.InfoKV(ctx, fmt.Sprintf("Status server is running on %s", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WarnKV(ctx, "Failed running status server", "err", err)
			cancel()
		}
	}()

	restServer, err := rest.NewRestServer()
	if err != nil {
		logger.FatalKV(ctx, "Сould not create a new rest server", "err", err)
	}

	restAddr := fmt.Sprintf("%s:%v",
		cfg.Services.Rest.Host,
		cfg.Services.Rest.Port,
	)

	srv := &http.Server{
		Addr:    restAddr,
		Handler: restServer.Router,
	}

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Rest Server is listening on: %s", restAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalKV(ctx, "Failed running rest server", "err", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.InfoKV(ctx, "The service is ready to accept requests")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}

	logger.InfoKV(ctx, "Shutting down server...")

	isReady.Store(false)

	if err := srv.Shutdown(ctx); err != nil {
		logger.FatalKV(ctx, "restServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "Rest server shut down correctly")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "statusServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "Status server shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.FatalKV(ctx, "metricsServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "Metrics server shut down correctly")
	}

}
