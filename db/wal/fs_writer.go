package wal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

var now = time.Now

type FSWriter struct {
	segment        *os.File
	directory      string
	segmentSize    int
	maxSegmentSize int
	log            *slog.Logger
}

func NewFSWriter(directory string, maxSegmentSize int, log *slog.Logger) *FSWriter {
	return &FSWriter{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
		log:            log,
	}
}

func (w *FSWriter) WriteBatch(batch []Log) {
	if w.segment == nil {
		if err := w.rotateSegment(); err != nil {
			w.acknowledgeWrite(batch, err)
			return
		}
	}

	if w.segmentSize >= w.maxSegmentSize {
		if err := w.rotateSegment(); err != nil {
			w.acknowledgeWrite(batch, err)
			return
		}
	}

	logs := make([]models.LogData, 0, len(batch))

	for i := 0; i < len(batch); i++ {
		logs = append(logs, batch[i].data)
	}

	if err := w.writeLogs(logs); err != nil {
		w.acknowledgeWrite(batch, err)
		return
	}

	err := w.segment.Sync() //системный вызов, который доноси данные на диск
	if err != nil {
		w.log.Error("failed to sync segment file", "err", err)
	}

	w.acknowledgeWrite(batch, err)
}

func (w *FSWriter) writeLogs(logs []models.LogData) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(logs); err != nil {
		w.log.Error("failed to write logs data", "err", err)
		return err
	}

	writtenBytes, err := w.segment.Write(buffer.Bytes())
	if err != nil {
		w.log.Error("failed to write logs data", "err", err)
		return err
	}
	w.segmentSize += writtenBytes
	return nil
}

func (w *FSWriter) acknowledgeWrite(batch []Log, err error) {
	//как только подвердили запись - отпускаем клиентов
	for _, log := range batch {
		log.SetResult(err)
	}
}

func (w *FSWriter) rotateSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", w.directory, now().UnixMilli())

	err := os.MkdirAll(w.directory, 0755)
	if err != nil {
		w.log.Error("failed to create directory for wal segment", "err", err)
		return err
	}

	flags := os.O_CREATE | os.O_WRONLY
	segment, err := os.OpenFile(segmentName, flags, 0644)
	if err != nil {
		w.log.Error("failed to create wal segment", "err", err)
		return err
	}

	w.segment = segment
	w.segmentSize = 0
	return nil
}
