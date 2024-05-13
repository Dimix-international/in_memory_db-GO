package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/db"
	"github.com/Dimix-international/in_memory_db-GO/internal/app"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/service"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)
	app := app.NewApp(
		log,
		handler.NewHanlderMessages(log, service.NewParserService(), service.NewAnalyzerService(), db.NewShardMap(10)),
	)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil {
			log.Error("Stop app", "err", err)
			exit <- syscall.SIGTERM
			close(exit)
		}
	}()

	<-exit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownApp(shutdownCtx); err != nil {
		log.Error("error closing app", "err", err)
	}

	log.Info("program finish")
}
