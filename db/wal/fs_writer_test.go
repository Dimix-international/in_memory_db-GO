package wal

import (
	"os"
	"testing"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/models"

	"github.com/stretchr/testify/assert"
)

const testWALDirectory = "temp_test_data"

func TestWriteBatch(t *testing.T) {
	writer := NewFSWriter(testWALDirectory, 10, logger.SetupLogger(""))
	batch := []Log{
		NewLog(int64(1), models.SetCommandID, []string{"key_4", "value_4"}),
		NewLog(int64(2), models.SetCommandID, []string{"key_2", "value_2"}),
		NewLog(int64(3), models.SetCommandID, []string{"key_3", "value_3"}),
	}

	go func() {
		writer.WriteBatch(batch)
	}()

	time.Sleep(time.Second * 3)

	for _, record := range batch {
		future := record.Result()
		err := future.Get()
		assert.Equal(t, nil, err)
	}

	files, err := os.ReadDir(testWALDirectory)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(files))
}
