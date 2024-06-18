package tools

import (
	"testing"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"

	"github.com/stretchr/testify/assert"
)

type testParserSizeItem struct {
	size   string
	result int
	err    error
}

func TestParseSize(t *testing.T) {
	testCases := [...]testParserSizeItem{
		{
			size:   "incorrect size",
			result: 0,
			err:    models.ErrIncorrectParseSize,
		},
		{
			size:   "20B",
			result: 20,
			err:    nil,
		},
		{
			size:   "20b",
			result: 20,
			err:    nil,
		},
		{
			size:   "20",
			result: 20,
			err:    nil,
		},
		{
			size:   "20KB",
			result: 20 * 1024,
			err:    nil,
		},
		{
			size:   "20Mb",
			result: 20 * 1024 * 1024,
			err:    nil,
		},
		{
			size:   "20gb",
			result: 20 * 1024 * 1024 * 1024,
			err:    nil,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run("check parse size", func(t *testing.T) {
			t.Parallel()

			size, err := ParseSize(test.size)

			assert.Equal(t, size, test.result)
			assert.Equal(t, err, test.err)
		})
	}
}
