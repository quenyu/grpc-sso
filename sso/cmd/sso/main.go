package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/quenyu/grpc-sso/sso/internal/app"
	"github.com/quenyu/grpc-sso/sso/internal/config"
	"github.com/quenyu/grpc-sso/sso/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := setupLogger(cfg.Env)
	logger.Info("starting sso server", "config", cfg)

	application := app.New(logger, cfg.Grpc.Port, cfg.StoragePath[0], cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.GRPCServer.Stop()
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
