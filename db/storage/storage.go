package storage

import "log/slog"

type db interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	Delete(key string)
}

type Storage struct {
	db  db
	log *slog.Logger
}

func NewStorage(db db, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
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
