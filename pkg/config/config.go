package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	GRPCPort string
	DBConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
}

func LoadConfig(configPath string) (*Config, error) {
	if err := godotenv.Load(configPath); err != nil {
		logrus.WithError(err).Warn("failed to load config file, using environment variables")
	}

	cfg := &Config{
		GRPCPort: os.Getenv("GRPC_PORT"),
		DBConfig: struct {
			Host     string
			Port     string
			User     string
			Password string
			DBName   string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}

	if cfg.GRPCPort == "" {
		cfg.GRPCPort = ":50051"
		logrus.Warn("GRPC_PORT not set, defaulting to :50051")
	}
	if cfg.DBConfig.Host == "" || cfg.DBConfig.User == "" || cfg.DBConfig.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}
