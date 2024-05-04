package main

import (
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/parser"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)

	pars := parser.NewParser(log)
	tokens, err := pars.Parse("DEL user_****")
	if err != nil {
		log.Error("parse error", "parser", err)
	}
	log.Info("app finish", "tokens:", tokens)
	log.Info("app finish")
}
