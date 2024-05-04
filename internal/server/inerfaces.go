package server

import "github.com/Dimix-international/in_memory_db-GO/internal/models"

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
