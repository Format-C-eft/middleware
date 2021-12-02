package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/Format-C-eft/middleware/internal/database/cache"
	"github.com/Format-C-eft/middleware/internal/logger"
	"github.com/Format-C-eft/middleware/internal/servers"
	"github.com/Format-C-eft/middleware/internal/tracer"
)

var configFile = flag.String("config", "", "Usage: -config=<config_file>")
var info = flag.Bool("version", false, "Usage: -info")

func main() {

	flag.Parse()

	if *info {
		fmt.Print(config.GetInfo())
		os.Exit(0)
	}

	ctx := context.Background()
	defer ctx.Done()

	err := config.ReadConfigYML(*configFile)

	if err != nil {
		logger.FatalKV(ctx, "Launch failed, settings file not read", "err", err)
	}

	cfg := config.GetConfigInstance()

	syncLogger := logger.InitLogger(&cfg)
	defer syncLogger()

	err = cache.InitClient(&cfg)
	if err != nil {
		logger.FatalKV(ctx, "Launch failed, cache DB not connect", "err", err)
	}

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		logger.FatalKV(ctx, "Could not create a new tracer", "err", err)
	}
	defer tracing.Close()

	servers.Start(ctx)
}
