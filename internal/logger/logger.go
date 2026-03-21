package logger

import "go.uber.org/zap"

func MustLoadLogger(env string) *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)

	switch env {
	case "local":
		logger = zap.NewExample()
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to init logger")
	}
	return logger
}
