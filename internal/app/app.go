package app

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/handler"
	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

// App - server instance for handling requests
type App struct {
	log            *slog.Logger
	handlerMessage *handler.HandlerMessages
	closers        []models.CloseFunc
}

// NewApp creates an instance of server
func NewApp(log *slog.Logger, handlerMessage *handler.HandlerMessages) *App {
	return &App{
		log:            log,
		handlerMessage: handlerMessage,
	}
}

// Run launches app
func (s *App) Run() {
	s.log = s.log.With(slog.String("op", "app"))
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.launchApp(); err != nil {
			s.log.Error("Stop server", "err", err)
			exit <- syscall.SIGTERM
			close(exit)
		}
	}()

	<-exit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.shutdown(shutdownCtx); err != nil {
		s.log.Error("error closing app", "err", err)
	}
}

func (s *App) launchApp() error {
	reader := bufio.NewReader(os.Stdin)
	s.log.Info("start app")

	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			s.log.Error("error read string", "err", err)
			continue
		}

		request = strings.TrimSpace(request)
		if len(request) == 0 {
			continue
		}

		s.handlerMessage.ProcessMessage(request)
	}
}

func (s *App) shutdown(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}
