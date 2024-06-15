package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/network"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Logging.Level)

	server, err := network.NewTCPServer(cfg, log)
	if err != nil {
		log.Error("error start VENOM", "error", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info("start VENOM")

	if err := server.Run(ctx); err != nil {
		log.Error("finish VENOM", "error", err)
	}

	log.Info("finish VENOM")
}
