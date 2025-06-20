package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/rdimidov/kvstore/internal/application/config"
	"github.com/rdimidov/kvstore/internal/application/services"

	"github.com/rdimidov/kvstore/internal/infrastructure/storage"
	"github.com/rdimidov/kvstore/internal/infrastructure/wal"
	"github.com/rdimidov/kvstore/internal/presentation/interpreter"
	"github.com/rdimidov/kvstore/internal/presentation/tcpserver"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := mustLoadConfig()
	defer cfg.Cleanup()

	logger := cfg.Logger()

	handler := mustInitHandler(ctx, cfg, logger)

	server := mustInitServer(cfg, handler)
	server.Start(ctx)
	logger.Infow("server exited gracefully")
}

func mustLoadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func mustInitHandler(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) *interpreter.RawInterpreter {
	repo := storage.NewMemory()

	var w services.WALogger
	if cfg.WAL.Enabled {

		walHandler, err := interpreter.New(repo)
		if err != nil {
			logger.Fatalw("failed to initialize interpreter", "error", err)
		}

		w, err = wal.New(ctx, cfg, walHandler)
		if err != nil {
			logger.Fatalw("failed to initialize wall", "error", err)
		}
	} else {
		w = &wal.Noop{}
	}

	app, err := services.NewApplication(ctx, repo, logger, w)
	if err != nil {
		logger.Fatalw("failed to initialize app", "error", err)
	}

	handler, err := interpreter.NewRaw(app)
	if err != nil {
		logger.Fatalw("failed to initialize interpreter", "error", err)
	}
	return handler
}

func mustInitServer(config *config.Config, handler *interpreter.RawInterpreter) *tcpserver.Server {
	logger := config.Logger()
	server, err := tcpserver.New(
		config.Network.Address, handler, logger,
		tcpserver.WithBufferSize(config.Network.MaxMessageSize),
		tcpserver.WithTimeouts(config.Network.ReadTimeout, config.Network.WriteTimeout),
	)
	if err != nil {
		logger.Fatalw("failed to create TCP server", "error", err)
	}
	return server
}
