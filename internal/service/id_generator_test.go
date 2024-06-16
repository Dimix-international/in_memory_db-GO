package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator()
	assert.Equal(t, int64(0), generator.counter.Load())

	id1 := generator.Generate()
	assert.Equal(t, int64(1), id1)

	id2 := generator.Generate()
	assert.Equal(t, int64(2), id2)
}

func TestSetInitID(t *testing.T) {
	t.Parallel()

	generator := NewIDGenerator()
	generator.SetInitValue(10)

	assert.Equal(t, int64(10), generator.counter.Load())

	generator.SetInitValue(50)
	assert.Equal(t, int64(10), generator.counter.Load())
}
