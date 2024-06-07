package main

import (
	"github.com/Dimix-international/in_memory_db-GO/db"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/network"
	"github.com/Dimix-international/in_memory_db-GO/internal/service"
)

const (
	shardValue = 10
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Logging.Level)
	handlerRequest := handler.NewHanlderMessages(log, service.NewParserService(), service.NewAnalyzerService(), db.NewShardMap(shardValue))

	server, err := network.NewTCPServer(handlerRequest, cfg.Network, log)
	if err != nil {
		log.Error("error start VENOM", "error", err)
	}

	if err := server.Run(); err != nil {
		log.Error("finish VENOM", "error", err)
	}

	log.Info("finish VENOM")
}
