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

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

// Server - server instance for handling requests
type Server struct {
	log      *slog.Logger
	parser   parserService
	analyzer analyzerService
	db       dbStorage
	closers  []models.CloseFunc
}

// NewServer creates an instance of server
func NewServer(log *slog.Logger, parser parserService, analyzer analyzerService, db dbStorage) *Server {
	return &Server{
		log:      log,
		parser:   parser,
		analyzer: analyzer,
		db:       db,
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

		s.handlerMessages(request)
	}
}

func (s *Server) handlerMessages(message string) {
	tokens, err := s.parser.Parse(message)
	if err != nil {
		s.log.Error("parsing error", "err", err)
		return
	}

	query, err := s.analyzer.Analyze(tokens)
	if err != nil {
		s.log.Error("analyzing error", "err", err)
		return
	}

	switch query.Command {
	case models.GetCommand:
		value, ok := s.db.Get(query.Arguments[0])
		if !ok {
			s.log.Info("key in db is not exist", "key", query.Arguments[0])
			break
		}
		s.log.Info("got value from db", "value", value)
	case models.SetCommand:
		s.db.Set(query.Arguments[0], query.Arguments[1])
		s.log.Info("command SET is execute", "key", query.Arguments[0])
	case models.DeleteCommand:
		s.db.Delete(query.Arguments[0])
		s.log.Info("command DELETE is execute", "key", query.Arguments[0])
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
