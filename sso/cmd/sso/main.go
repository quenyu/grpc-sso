package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/quenyu/grpc-sso/sso/internal"
	"github.com/quenyu/grpc-sso/sso/internal/lib/logger/handlers/slogpretty"
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
		log = setupPrettyLogger()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
