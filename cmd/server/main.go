package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krchnk/exchange-rates-service/internal/exchangerates"
	"github.com/Krchnk/exchange-rates-service/internal/storages/postgres"
	"github.com/Krchnk/exchange-rates-service/pkg/config"
	"github.com/Krchnk/exchange-rates-service/pkg/logger"
)

func main() {
	configPath := flag.String("c", "config.env", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.NewLogger().WithError(err).Fatal("failed to load config")
	}

	logger := logger.NewLogger()
	logger.WithField("config", cfg).Info("configuration loaded")

	store, err := postgres.NewStorage(cfg.DBConfig)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to database")
	}
	logger.Info("database connection established")

	service := exchangerates.NewService(store, logger)

	server, err := exchangerates.NewServer(cfg, service)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize gRPC server")
	}

	go func() {
		if err := server.Start(cfg.GRPCPort); err != nil {
			logger.WithError(err).Fatal("failed to start gRPC server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	server.Stop()
	logger.Info("Server stopped gracefully")
}
