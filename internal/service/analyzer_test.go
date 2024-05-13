package service

import (
	"fmt"
	"testing"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"

	"github.com/stretchr/testify/assert"
)

type testAnalyzeItem struct {
	tokens []string
	qwery  models.Query
	err    error
}

func TestAnalyze(t *testing.T) {
	testCases := [...]testAnalyzeItem{
		{
			tokens: []string{"SET", "weather_2_pm", "cold_moscow_weather"},
			qwery:  models.Query{Command: models.SetCommand, Arguments: []string{"weather_2_pm", "cold_moscow_weather"}},
			err:    nil,
		},
		{
			tokens: []string{"GET", "weather_2_pm"},
			qwery:  models.Query{Command: models.GetCommand, Arguments: []string{"weather_2_pm"}},
			err:    nil,
		},
		{
			tokens: []string{"DELETE", "weather_2_pm"},
			qwery:  models.Query{Command: models.DeleteCommand, Arguments: []string{"weather_2_pm"}},
			err:    nil,
		},
		{
			tokens: []string{},
			qwery:  models.Query{},
			err:    fmt.Errorf("arguments is empty: %w", models.ErrInvalidArguments),
		},
		{
			tokens: []string{"TRANSACTION", "weather_2_pm"},
			qwery:  models.Query{},
			err:    fmt.Errorf("command TRANSACTION is not exist: %w", models.ErrInvalidArguments),
		},
		{
			tokens: []string{"GET", "weather_2_pm", "cold_moscow_weather"},
			qwery:  models.Query{},
			err: fmt.Errorf(
				"passed number of arguments equal to 2 is not correct for command GET: %w",
				models.ErrInvalidArguments,
			),
		},
		{
			tokens: []string{"SET", "weather_2_pm", "cold_moscow_weather", "cold_moscow_weather cold"},
			qwery:  models.Query{},
			err: fmt.Errorf(
				"passed number of arguments equal to 3 is not correct for command SET: %w",
				models.ErrInvalidArguments,
			),
		},
		{
			tokens: []string{"DELETE"},
			qwery:  models.Query{},
			err: fmt.Errorf(
				"passed number of arguments equal to 0 is not correct for command DELETE: %w",
				models.ErrInvalidArguments,
			),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run("check analyzing", func(t *testing.T) {
			t.Parallel()

			analyzer := NewAnalyzerService()
			qwery, err := analyzer.Analyze(test.tokens)

			assert.Equal(t, qwery, test.qwery)
			assert.Equal(t, err, test.err)
		})
	}
}
