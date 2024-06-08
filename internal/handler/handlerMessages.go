package handler

import (
	"fmt"

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
	parser   parserService
	analyzer analyzerService
	db       dbStorage
}

// NewHanlderMessages creating handler instance
func NewHanlderMessages(parser parserService, analyzer analyzerService, db dbStorage) *HandlerMessages {
	return &HandlerMessages{
		parser:   parser,
		analyzer: analyzer,
		db:       db,
	}
}

// ProcessMessage start work with message
func (s *HandlerMessages) ProcessMessage(command []byte) string {
	tokens, err := s.parser.Parse(string(command))
	if err != nil {
		return fmt.Sprintf("parsing error: %v", err)
	}

	query, err := s.analyzer.Analyze(tokens)
	if err != nil {
		return fmt.Sprintf("analyzing error: %v", err)
	}

	switch query.Command {
	case models.GetCommand:
		value, ok := s.db.Get(query.Arguments[0])
		if !ok {
			return fmt.Sprintf("key in db is not exist: %v", query.Arguments[0])
		}
		return fmt.Sprintf("got value from db: %v", value)
	case models.SetCommand:
		s.db.Set(query.Arguments[0], query.Arguments[1])
		return fmt.Sprintf("command SET is execute: %v", query.Arguments[0])
	case models.DeleteCommand:
		s.db.Delete(query.Arguments[0])
		return fmt.Sprintf("command DELETE is execute: %v", query.Arguments[0])
	}

	return fmt.Sprintf("unknown command: %v", query.Arguments[0])
}
