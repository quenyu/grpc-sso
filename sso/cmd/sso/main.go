package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/quenyu/grpc-sso/sso/internal"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("error initializing config: %s", err)
	}

	logger := setupLogger(cfg.Env)
	logger.Info("starting sso server", "config", cfg)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return log
}
