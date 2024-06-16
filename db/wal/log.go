package wal

import (
	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/Dimix-international/in_memory_db-GO/internal/tools"
)

type Log struct {
	data         models.LogData
	writePromise tools.Promise
}

func NewLog(lsn int64, commandID int, arguments []string) Log {
	return Log{
		data: models.LogData{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: arguments,
		},
		writePromise: *tools.NewPromise(),
	}
}

func (l *Log) Data() models.LogData {
	return l.data
}

func (l *Log) CommandID() int {
	return l.data.CommandID
}

func (l *Log) LSN() int64 {
	return l.data.LSN
}

func (l *Log) Arguments() []string {
	return l.data.Arguments
}

func (l *Log) SetResult(err error) {
	l.writePromise.Set(err)
}

func (l *Log) Result() tools.Future {
	return *l.writePromise.GetFuture()
}
