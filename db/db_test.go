package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBMap(t *testing.T) {
	db := NewDBMap()

	db.Set("a", "1")
	db.Set("b", "2")
	db.Set("c", "3")

	value, ok := db.Get("a")
	assert.True(t, ok)
	assert.Equal(t, value, "1")

	db.Delete("b")
	value, ok = db.Get("b")

	assert.Equal(t, ok, false)
	assert.Equal(t, value, "")

	db.Set("c", "4")
	value, ok = db.Get("c")

	assert.True(t, ok)
	assert.Equal(t, value, "4")
}
