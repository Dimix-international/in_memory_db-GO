package handler

import (
	"context"
	"fmt"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

type parserService interface {
	Parse(query string) ([]string, error)
}

type analyzerService interface {
	Analyze(tokens []string) (models.Query, error)
}

type storage interface {
	Get(key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
}

// HandlerMessages - handler instance for handling messages
type HandlerMessages struct {
	parser   parserService
	analyzer analyzerService
	store    storage
}

// NewHanlderMessages creating handler instance
func NewHanlderMessages(parser parserService, analyzer analyzerService, store storage) *HandlerMessages {
	return &HandlerMessages{
		parser:   parser,
		analyzer: analyzer,
		store:    store,
	}
}

// ProcessMessage start work with message
func (s *HandlerMessages) ProcessMessage(ctx context.Context, command []byte) string {
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
		value, _ := s.store.Get(query.Arguments[0])
		if len(value) == 0 {
			return fmt.Sprintf("key in db is not exist: %v", query.Arguments[0])
		}
		return fmt.Sprintf("got value from db: %v", value)
	case models.SetCommand:
		err := s.store.Set(ctx, query.Arguments[0], query.Arguments[1])
		if err != nil {
			return fmt.Sprintf("falied SET command: %v with error %v", query.Arguments[0], err)
		}
		return fmt.Sprintf("command SET is execute: %v", query.Arguments[0])
	case models.DeleteCommand:
		err := s.store.Del(ctx, query.Arguments[0])
		if err != nil {
			return fmt.Sprintf("falied DELETE command: %v with error %v", query.Arguments[0], err)
		}
		return fmt.Sprintf("command DELETE is execute: %v", query.Arguments[0])
	}

	return fmt.Sprintf("unknown command: %v", query.Arguments[0])
}
