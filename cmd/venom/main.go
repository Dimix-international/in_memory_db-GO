package main

import (
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/network"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Logging.Level)

	server, err := network.NewTCPServer(cfg.Network, log)
	if err != nil {
		log.Error("error start VENOM", "error", err)
	}

	if err := server.Run(); err != nil {
		log.Error("finish VENOM", "error", err)
	}

	log.Info("finish VENOM")
}
