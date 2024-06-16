package wal

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

type fsWriter interface {
	WriteBatch([]Log)
}

type fsReader interface {
	ReadLogs() ([]models.LogData, error)
}

type WAL struct {
	fsWriter     fsWriter
	fsReader     fsReader
	flushTimeout time.Duration //по истеканию записываем в log
	maxBatchSize int

	mutex   sync.Mutex
	batch   []Log
	batches chan []Log

	closeCh     chan struct{}
	closeDoneCh chan struct{}

	log *slog.Logger
}

func NewWAL(
	fsWriter fsWriter,
	fsReader fsReader,
	flushTimeout time.Duration,
	maxBatchSize int,
	log *slog.Logger,
) *WAL {
	wall := &WAL{
		fsWriter:     fsWriter,
		fsReader:     fsReader,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []Log),
		closeCh:      make(chan struct{}),
		closeDoneCh:  make(chan struct{}),
		log:          log,
	}

	return wall
}

func (w *WAL) Start() {
	go func() {
		defer close(w.closeDoneCh)
		timer := time.NewTimer(w.flushTimeout)

		for {
			timer.Reset(w.flushTimeout)

			select {
			case <-w.closeCh:
				w.flushBatch()
				return
			case batch := <-w.batches:
				w.fsWriter.WriteBatch(batch)
			case <-timer.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Shutdown() {
	close(w.closeCh)
	<-w.closeDoneCh
}

func (w *WAL) flushBatch() {
	var batch []Log
	tools.WithLock(&w.mutex, func() {
		if len(w.batch) != 0 {
			batch = w.batch
			w.batch = nil
		}
	})

	if len(batch) != 0 {
		w.fsWriter.WriteBatch(batch)
	}
}

func (w *WAL) Set(ctx context.Context, key, value string) tools.Future {
	return w.push(ctx, models.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) tools.Future {
	return w.push(ctx, models.DeleteCommandID, []string{key})
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) tools.Future {
	txID := ctx.Value(models.KeyTxID).(int64)
	record := NewLog(txID, commandID, args)

	tools.WithLock(&w.mutex, func() {
		w.batch = append(w.batch, record)

		if len(w.batch) >= w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return record.Result()
}

func (w *WAL) TryRecoverWALSegments(stream chan<- []models.LogData) {
	logs, err := w.fsReader.ReadLogs()
	if err != nil {
		w.log.Error("failed to recover WAL segments", "err", err)
	} else {
		stream <- logs
	}
}
