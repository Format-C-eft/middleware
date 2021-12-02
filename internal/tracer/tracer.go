package tracer

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerZap "github.com/uber/jaeger-client-go/log/zap"
)

// NewTracer - new jaeger
func NewTracer(cfg *config.Config) (io.Closer, error) {

	cfgTracer := &jaegercfg.Configuration{
		Disabled:    !cfg.Servers.Jaeger.Use,
		ServiceName: cfg.Servers.Jaeger.Service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  cfg.Servers.Jaeger.Host + cfg.Servers.Jaeger.Port,
		},
	}

	jaegerLogger := jaegerZap.NewLogger(logger.GetNotSugaredLogger())
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaegerLogger))
	if err != nil {
		logger.FatalKV(context.TODO(), "failed init jaeger", "err", err)
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

// OperationSetName - for jaeger correct names
func OperationSetName(r *http.Request) string {

	cfg := config.GetConfigInstance()

	return strings.ToUpper(strings.Replace(r.URL.Path, cfg.Services.Rest.Path, "", 1))
}
