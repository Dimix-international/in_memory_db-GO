package storage

import (
	"context"
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
	Shutdown()
}

type idGenerator interface {
	SetInitValue(inintValue int64)
	Generate() int64
}

type Storage struct {
	db          db
	wal         wal
	idGenerator idGenerator
	log         *slog.Logger
}

func NewStorage(db db, wal wal, idGenerator idGenerator, log *slog.Logger) *Storage {
	return &Storage{
		db:          db,
		wal:         wal,
		idGenerator: idGenerator,
		log:         log,
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
	var (
		maxID    int64 = 0
		logsChan       = make(chan []models.LogData)
	)

	go func() {
		defer close(logsChan)
		s.wal.TryRecoverWALSegments(logsChan)
	}()

	for logs := range logsChan {
		maxID = logs[len(logs)-1].LSN

		for i := 0; i < len(logs); i++ {
			switch logs[i].CommandID {
			case models.SetCommandID:
				s.db.Set(logs[i].Arguments[0], logs[i].Arguments[1])
			case models.DeleteCommandID:
				s.db.Delete(logs[i].Arguments[0])
			}
		}
	}

	s.idGenerator.SetInitValue(maxID)
}

func (s *Storage) TransactionID() int64 {
	return s.idGenerator.Generate()
}

func (s *Storage) Shutdown() {
	s.wal.Shutdown()
}
