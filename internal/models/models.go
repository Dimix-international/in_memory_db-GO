package models

import "errors"

const (
	GetCommand     = "GET"
	SetCommand     = "SET"
	DeleteCommand  = "DELETE"
	UnknownCommand = "UNKNOWN"

	GetCommandArgumentsNumber    = 1
	DeleteCommandArgumentsNumber = 1
	SetCommandArgumentsNumber    = 2
)

var (
	CommandRatioWithArument = map[string]int{
		GetCommand:    GetCommandArgumentsNumber,
		SetCommand:    SetCommandArgumentsNumber,
		DeleteCommand: DeleteCommandArgumentsNumber,
	}
	ErrParsing          = errors.New("invalid argument for parsing")
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrInvalidCommand   = errors.New("invalid command")
)
