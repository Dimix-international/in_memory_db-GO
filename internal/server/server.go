package server

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

// Server - server instance for handling requests
type Server struct {
	log            *slog.Logger
	handlerMessage *handler.HandlerMessages
	closers        []models.CloseFunc
}

// NewServer creates an instance of server
func NewServer(log *slog.Logger, handlerMessage *handler.HandlerMessages) *Server {
	return &Server{
		log:            log,
		handlerMessage: handlerMessage,
	}
}

// Run launches server
func (s *Server) Run() {
	s.log = s.log.With(slog.String("op", "server"))
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.launchServer(); err != nil {
			s.log.Error("Stop server", "err", err)
			exit <- syscall.SIGTERM
			close(exit)
		}
	}()

	<-exit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.shutdown(shutdownCtx); err != nil {
		s.log.Error("error closing server", "err", err)
	}
}

func (s *Server) launchServer() error {
	reader := bufio.NewReader(os.Stdin)
	s.log.Info("start server")

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

func (s *Server) shutdown(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}
