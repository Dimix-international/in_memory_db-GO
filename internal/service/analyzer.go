package service

import (
	"fmt"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

// AnalyzerService - service for analyzing messages
type AnalyzerService struct{}

// NewAnalyzerService creates an instance of the message analyser
func NewAnalyzerService() *AnalyzerService {
	return &AnalyzerService{}
}

// Analyze launching the message analysis method
func (a *AnalyzerService) Analyze(tokens []string) (models.Query, error) {
	if len(tokens) == 0 {
		return models.Query{}, fmt.Errorf("arguments is empty: %w", models.ErrInvalidArguments)
	}

	countArguments, ok := models.CommandRatioWithArument[tokens[0]]

	if !ok {
		return models.Query{}, fmt.Errorf("command %v is not exist: %w", tokens[0], models.ErrInvalidArguments)
	}

	if len(tokens[1:]) != countArguments {
		return models.Query{}, fmt.Errorf(
			"passed number of arguments equal to %d is not correct for command %v: %w",
			len(tokens[1:]),
			tokens[0],
			models.ErrInvalidArguments,
		)
	}

	return models.Query{Command: tokens[0], Arguments: tokens[1:]}, nil
}
