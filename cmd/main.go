package main

import (
	"github.com/Dimix-international/in_memory_db-GO/db"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/server"
	"github.com/Dimix-international/in_memory_db-GO/internal/service"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)
	server.NewServer(
		log,
		handler.NewHanlderMessages(log, service.NewParserService(), service.NewAnalyzerService(log), db.NewShardMap(10)),
	).Run()

	log.Info("app finish")
}
