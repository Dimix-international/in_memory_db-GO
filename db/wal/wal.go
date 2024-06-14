package wal

import (
	"log/slog"
	"sync"
	"time"
)

type fsWriter interface {
	WriteBatch([]Log)
}

type fsReader interface {
	ReadLogs() ([]LogData, error)
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
	return &WAL{
		fsWriter:     fsWriter,
		fsReader:     fsReader,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []Log),
		closeCh:      make(chan struct{}),
		closeDoneCh:  make(chan struct{}),
		log:          log,
	}
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

func (w *WAL) flushBatch() {}

func (w *WAL) tryRecoverWALSegments(stream chan<- []LogData) {
	logs, err := w.fsReader.ReadLogs()
	if err != nil {
		w.log.Error("failed to recover WAL segments", "err", err)
	} else {
		stream <- logs
	}
}
