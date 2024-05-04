package parser

import (
	"log/slog"
	"testing"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
	"github.com/stretchr/testify/assert"
)

type TestItems struct {
	query  string
	tokens []string
	err    error
}

func Test_Parse(t *testing.T) {
	testCases := [...]TestItems{
		{query: "SET weather_2_pm cold_moscow_weather", tokens: []string{"SET", "weather_2_pm", "cold_moscow_weather"}, err: nil},
		{query: "SET\tweather_2_pm\tcold_moscow_weather", tokens: []string{"SET", "weather_2_pm", "cold_moscow_weather"}, err: nil},
		{query: "GET /etc/nginx/config.yaml", tokens: nil, err: models.ErrParsing},
		{query: "GET ++", tokens: nil, err: models.ErrParsing},
		{query: "GET   ", tokens: []string{"GET"}, err: nil},
		{query: "   GET ++", tokens: nil, err: models.ErrParsing},
		{query: "DEL user_****", tokens: []string{"DEL", "user_****"}, err: nil},
		{query: "DEL\nuser_****\nmany_words\n", tokens: []string{"DEL", "user_****", "many_words"}, err: nil},
		{query: "SET user_**** cre++ate", tokens: nil, err: models.ErrParsing},
	}

	for _, test := range testCases {
		test := test
		t.Run("check correcting of parsing", func(t *testing.T) {
			t.Parallel()

			parser := NewParser(slog.Default())
			tokens, err := parser.Parse(test.query)

			assert.Equal(t, tokens, test.tokens)
			assert.Equal(t, err, test.err)
		})
	}
}
