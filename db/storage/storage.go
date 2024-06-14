package storage

import "log/slog"

type db interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	Delete(key string)
}

type wal interface {
	Start()
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
		s.wal.Start()
	}
}

func (s *Storage) Get(key string) (string, error) {
	value, _ := s.db.Get(key)
	return value, nil
}

func (s *Storage) Set(key, value string) error {
	s.db.Set(key, value)
	return nil
}

func (s *Storage) Del(key string) error {
	s.db.Delete(key)
	return nil
}
