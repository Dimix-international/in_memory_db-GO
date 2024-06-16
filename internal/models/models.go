package models

import (
	"context"
	"errors"
)

const (
	// GetCommand the command to get data
	GetCommand = "GET"
	// SetCommand the command to set data
	SetCommand = "SET"
	// DeleteCommand the command to delete data
	DeleteCommand = "DELETE"
	// GetCommandArgumentsNumber the number of arguments for the GET command
	GetCommandArgumentsNumber = 1
	// DeleteCommandArgumentsNumber the number of arguments for the DELETE command
	DeleteCommandArgumentsNumber = 1
	// SetCommandArgumentsNumber the number of arguments for the SET command
	SetCommandArgumentsNumber = 2
)

var (
	// CommandRatioWithArument map for checking existing of command
	CommandRatioWithArument = map[string]int{
		GetCommand:    GetCommandArgumentsNumber,
		SetCommand:    SetCommandArgumentsNumber,
		DeleteCommand: DeleteCommandArgumentsNumber,
	}
	// ErrParsing error parsing the message
	ErrParsing = errors.New("invalid argument for parsing")
	// ErrInvalidArguments error for invalid arguments in the executing command to the database
	ErrInvalidArguments = errors.New("invalid arguments")
	// ErrInvalidCommand invalid command to write to the database
	ErrInvalidCommand = errors.New("invalid command")
	// ErrInvalidMaxConnections invalid count connections
	ErrInvalidMaxConnections = errors.New("invalid number of max connections")
	// ErrInvalidLogger invalid logger
	ErrInvalidLogger = errors.New("logger is invalid")
	// ErrInvalidLogger invalid logger
	ErrNetClosed = errors.New("net closed")
)

// Query - the structure that stores the command and its arguments for writing to the database
type Query struct {
	Command   string
	Arguments []string
}

// LogData - structure of logs for WAL
type LogData struct {
	LSN         int64 //id transaction
	CommandName string
	Arguments   []string
}

type txContextKey string

var KeyTxID = txContextKey("tx")

// CloseFunc - function for graceful shutdown
type CloseFunc func(context.Context) error
