package handler

import (
	"fmt"
	"log/slog"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

type parserService interface {
	Parse(query string) ([]string, error)
}

type analyzerService interface {
	Analyze(tokens []string) (models.Query, error)
}

type dbStorage interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	Delete(key string)
}

// HandlerMessages - handler instance for handling messages
type HandlerMessages struct {
	log      *slog.Logger
	parser   parserService
	analyzer analyzerService
	db       dbStorage
}

// NewHanlderMessages creating handler instance
func NewHanlderMessages(log *slog.Logger, parser parserService, analyzer analyzerService, db dbStorage) *HandlerMessages {
	return &HandlerMessages{
		log:      log,
		parser:   parser,
		analyzer: analyzer,
		db:       db,
	}
}

// ProcessMessage start work with message
func (s *HandlerMessages) ProcessMessage(message string) string {
	tokens, err := s.parser.Parse(message)
	if err != nil {
		s.log.Error("parsing error", "err", err)
		return fmt.Sprintf("parsing error: %v", err)
	}

	query, err := s.analyzer.Analyze(tokens)
	if err != nil {
		s.log.Error("analyzing error", "err", err)
		return fmt.Sprintf("analyzing error: %v", err)
	}

	switch query.Command {
	case models.GetCommand:
		value, ok := s.db.Get(query.Arguments[0])
		if !ok {
			s.log.Info("key in db is not exist", "key", query.Arguments[0])
			return fmt.Sprintf("key in db is not exist: %v", query.Arguments[0])
		}
		s.log.Info("got value from db", "value", value)
		return fmt.Sprintf("got value from db: %v", value)
	case models.SetCommand:
		s.db.Set(query.Arguments[0], query.Arguments[1])
		s.log.Info("command SET is execute", "key", query.Arguments[0])
		return fmt.Sprintf("command SET is execute: %v", query.Arguments[0])
	case models.DeleteCommand:
		s.db.Delete(query.Arguments[0])
		s.log.Info("command DELETE is execute", "key", query.Arguments[0])
		return fmt.Sprintf("command DELETE is execute: %v", query.Arguments[0])
	}

	s.log.Info("unknown command", "key", query.Arguments[0])
	return "unknown "
}
