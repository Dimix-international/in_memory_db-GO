package wal

import (
	"os"
	"testing"
	"time"

	"github.com/Dimix-international/in_memory_db-GO/internal/logger"
	"github.com/Dimix-international/in_memory_db-GO/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestFSReader(t *testing.T) {
	os.RemoveAll(testWALDirectory)

	writer := NewFSWriter(testWALDirectory, 10, logger.SetupLogger(""))
	batch := []Log{
		NewLog(int64(1), models.SetCommandID, []string{"key_1", "value_1"}),
		NewLog(int64(2), models.DeleteCommandID, []string{"key_2"}),
		NewLog(int64(3), models.SetCommandID, []string{"key_3", "value_3"}),
		NewLog(int64(4), models.DeleteCommandID, []string{"key_4"}),
	}

	go func() {
		writer.WriteBatch(batch)
	}()

	time.Sleep(time.Second * 3)

	reader := NewFSReader(testWALDirectory, logger.SetupLogger(""))

	logs, err := reader.ReadLogs()

	assert.NoError(t, err)
	assert.Equal(t, 4, len(logs))
	assert.Equal(t, "key_1", logs[0].Arguments[0])
	assert.Equal(t, models.DeleteCommandID, logs[1].CommandID)
	assert.Equal(t, "value_3", logs[2].Arguments[1])
	assert.Equal(t, int64(4), logs[3].LSN)
}
