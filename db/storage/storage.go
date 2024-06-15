package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

type db interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	Delete(key string)
}

type wal interface {
	Start()
	Set(ctx context.Context, key, value string) tools.Future
	Del(ctx context.Context, key string) tools.Future
	TryRecoverWALSegments(stream chan<- []models.LogData)
}

type Storage struct {
	db  db
	wal wal
	log *slog.Logger
}

func NewStorage(db db, wal wal, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		wal: wal,
		log: log,
	}
}

func (s *Storage) Start() {
	if s.wal != nil {
		s.recoverDB()
		s.wal.Start()
	}
}

func (s *Storage) Get(key string) (string, error) {
	value, _ := s.db.Get(key)
	return value, nil
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.wal != nil {
		future := s.wal.Set(ctx, key, value)
		if err := future.Get(); err != nil {
			return err.(error)
		}
	}

	s.db.Set(key, value)
	return nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if s.wal != nil {
		future := s.wal.Del(ctx, key)
		if err := future.Get(); err != nil {
			return err.(error)
		}
	}

	s.db.Delete(key)
	return nil
}

func (s *Storage) recoverDB() {
	logsChan := make(chan []models.LogData)

	go func() {
		defer close(logsChan)
		s.wal.TryRecoverWALSegments(logsChan)
	}()

	for logs := range logsChan {
		fmt.Println("logs", logs)
	}
}
