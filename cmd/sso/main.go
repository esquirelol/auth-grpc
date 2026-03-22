package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/esquirelol/auth-grpc/internal/app"
	"github.com/esquirelol/auth-grpc/internal/config"
	log "github.com/esquirelol/auth-grpc/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	logger := log.MustLoadLogger(cfg.Env)

	application := app.New(ctx, logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
}
