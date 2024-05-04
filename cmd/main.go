package main

import (
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)

	log.Info("app finish")
}
