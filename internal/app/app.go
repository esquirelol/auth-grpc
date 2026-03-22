package app

import (
	"context"
	"time"

	grpcapp "github.com/esquirelol/auth-grpc/internal/app/grpc"
	"github.com/esquirelol/auth-grpc/internal/services/auth"
	storage2 "github.com/esquirelol/auth-grpc/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(ctx context.Context, log *zap.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := storage2.New(ctx, storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, &storage, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
