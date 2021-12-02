package rest

import (
	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
)

type restServer struct {
	Router *gin.Engine
}

// NewRestServer - new rest server
func NewRestServer() (*restServer, error) {

	cfg := config.GetConfigInstance()
	if !cfg.Project.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginRecovery(true))
	router.Use(ginLogger())
	router.Use(
		ginhttp.Middleware(
			opentracing.GlobalTracer(),
			ginhttp.OperationNameFunc(tracer.OperationSetName),
		),
	)
	router.Use(SetHeaders())
	router.Use(AbortMetodOption())
	router.Use(AddSessionInfo())
	router.Use(ClearPath())
	router.Use(IncMetrics())

	initRoutes(router, &cfg)

	return &restServer{
		Router: router,
	}, nil
}
