package wal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

type FSReader struct {
	directory string
	log       *slog.Logger
}

func NewFSReader(directory string, log *slog.Logger) *FSReader {
	return &FSReader{
		directory: directory,
		log:       log,
	}
}

func (r *FSReader) ReadLogs() ([]models.LogData, error) {
	files, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	var logs []models.LogData

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		fileName := fmt.Sprintf("%s/%s", r.directory, files[i].Name())

		segmentLogs, err := r.readSegment(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to recove WAL segment: %w", err)
		}

		logs = append(logs, segmentLogs...)
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].LSN < logs[j].LSN
	})

	return logs, nil
}

func (r *FSReader) readSegment(filename string) ([]models.LogData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var logs []models.LogData
	buffer := bytes.NewBuffer(data)

	for buffer.Len() > 0 {
		var batch []models.LogData

		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&batch); err != nil {
			return nil, fmt.Errorf("failed to parse logs data: %w", err)
		}

		logs = append(logs, batch...)
	}

	return logs, nil
}
