package grpc

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/quenyu/grpc-sso/sso/internal/grpc/auth-service"
	authservice "github.com/quenyu/grpc-sso/sso/internal/grpc/auth-service"
	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	grpc *grpc.Server
	port int
}

func NewApp(
	log *slog.Logger,
	authservice authservice.Authservice,
	port int,
) *App {
	grpcServer := grpc.NewServer()

	authgrpc.Register(grpcServer, authservice)

	return &App{
		log:  log,
		grpc: grpcServer,
		port: port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.grpc.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.grpc.GracefulStop()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
