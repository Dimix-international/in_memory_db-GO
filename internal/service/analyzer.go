package service

import (
	"log/slog"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

type AnalyzerService struct {
	log *slog.Logger
}

func NewAnalyzerService(log *slog.Logger) *AnalyzerService {
	return &AnalyzerService{log: log}
}

func (a *AnalyzerService) Analyze(tokens []string) (models.Query, error) {
	a.log = a.log.With(slog.String("op", "service.Analyze"))

	if len(tokens) == 0 {
		a.log.Error("arguments is empy")
		return models.Query{}, models.ErrInvalidArguments
	}

	countArguments, ok := models.CommandRatioWithArument[tokens[0]]

	if !ok {
		a.log.Error("command is not exist")
		return models.Query{}, models.ErrInvalidCommand
	}

	if len(tokens[1:]) != countArguments {
		a.log.Error("count of arguments is not correct")
		return models.Query{}, models.ErrInvalidArguments
	}

	return models.Query{Command: tokens[0], Arguments: tokens[1:]}, nil
}
