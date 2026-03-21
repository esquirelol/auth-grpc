package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/numbergroup/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" end-default:"local"`
	StoragePath string
	TokenTTl    time.Duration `yaml:"token_ttl" env-default:"24h"`
	GRPC        ConfigGRPC    `yaml:"grpc"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port" env-default:"9091"`
	TimeOut time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic("")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic("file is not exists")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config")
	}

	cfg.StoragePath = os.Getenv("STORAGE_PATH")

	return &cfg
}
