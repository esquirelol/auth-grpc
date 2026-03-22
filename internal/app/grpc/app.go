package grpcapp

import (
	"fmt"
	"net"

	authrpc "github.com/esquirelol/auth-grpc/internal/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, port int, authService authrpc.Auth) *App {
	gRPCServer := grpc.NewServer()
	authrpc.RegisterServerAPI(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("")
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(zap.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("error:", zap.Error(err))
		return err
	}

	log.Info("grpc server is running", zap.String("addr", l.Addr().String()))
	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("error:", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
