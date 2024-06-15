package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/db/wal"
	"github.com/Dimix-international/in_memory_db-GO/internal/config"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./data/venom/wal"
)

func CreateWal(cfg *config.WALConfig, log *slog.Logger) (*wal.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg != nil {
		if cfg.FlushingBatchSize != 0 {
			flushingBatchSize = cfg.FlushingBatchSize
		}

		if cfg.FlushingBatchTimeout != 0 {
			flushingBatchTimeout = cfg.FlushingBatchTimeout
		}

		if cfg.MaxSegmentSize != "" {
			size, err := tools.ParseSize(cfg.MaxSegmentSize)
			if err != nil {
				return nil, errors.New("max segment size is incorrect")
			}

			maxSegmentSize = size
		}

		if cfg.DataDirectory != "" {
			dataDirectory = cfg.DataDirectory
		}

		return wal.NewWAL(
			wal.NewFSWriter(dataDirectory, maxSegmentSize, log),
			wal.NewFSReader(dataDirectory, log),
			flushingBatchTimeout,
			flushingBatchSize,
			log,
		), nil
	}

	return nil, nil
}
