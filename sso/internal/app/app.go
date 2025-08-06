package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/quenyu/grpc-sso/sso/internal/app/grpc"
	"github.com/quenyu/grpc-sso/sso/internal/services/auth"
	"github.com/quenyu/grpc-sso/sso/internal/storage/postgresql"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgresql.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
