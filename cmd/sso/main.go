package main

import (
	"github.com/esquirelol/auth-grpc/internal/config"
	log "github.com/esquirelol/auth-grpc/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	logger := log.MustLoadLogger(cfg.Env)
	logger.Info("logger init")
}
