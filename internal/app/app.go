package app

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"strings"

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

// Run app
func (s *App) Run() error {
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

func (s *App) ShutdownApp(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}
